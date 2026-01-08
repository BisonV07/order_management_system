package controllers

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"oms/server/api/v1/helpers"
	"gorm.io/gorm"
)

// MetricsController handles system metrics (Docker and PostgreSQL)
type MetricsController struct {
	db *gorm.DB
}

// NewMetricsController creates a new MetricsController
func NewMetricsController(db *gorm.DB) *MetricsController {
	return &MetricsController{db: db}
}

// GetDockerMetrics handles GET /api/v1/admin/metrics/docker - Get Docker container metrics
func (mc *MetricsController) GetDockerMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" && role != "ADMIN" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	// Get Docker stats for PostgreSQL container
	containerName := "oms_postgres"
	metrics := map[string]interface{}{
		"container_name": containerName,
		"timestamp":      time.Now().Unix(),
	}

	// Get container stats
	cmd := exec.CommandContext(ctx, "docker", "stats", containerName, "--no-stream", "--format", "{{.CPUPerc}},{{.MemUsage}},{{.MemPerc}},{{.NetIO}},{{.BlockIO}}")
	output, err := cmd.Output()
	if err != nil {
		// Container might not be running or docker command failed
		metrics["error"] = fmt.Sprintf("Failed to get Docker stats: %v", err)
		metrics["status"] = "unavailable"
		helpers.WriteJSONResponse(w, http.StatusOK, metrics)
		return
	}

	// Parse stats output
	stats := strings.TrimSpace(string(output))
	if stats != "" {
		parts := strings.Split(stats, ",")
		if len(parts) >= 5 {
			metrics["cpu_percent"] = strings.TrimSpace(parts[0])
			metrics["memory_usage"] = strings.TrimSpace(parts[1])
			metrics["memory_percent"] = strings.TrimSpace(parts[2])
			metrics["network_io"] = strings.TrimSpace(parts[3])
			metrics["block_io"] = strings.TrimSpace(parts[4])
		}
	}

	// Get container info
	infoCmd := exec.CommandContext(ctx, "docker", "inspect", containerName, "--format", "{{.State.Status}},{{.State.StartedAt}},{{.HostConfig.Memory}},{{.HostConfig.CpuShares}}")
	infoOutput, err := infoCmd.Output()
	if err == nil && len(infoOutput) > 0 {
		infoParts := strings.Split(strings.TrimSpace(string(infoOutput)), ",")
		if len(infoParts) >= 2 {
			metrics["status"] = infoParts[0]
			metrics["started_at"] = infoParts[1]
			if len(infoParts) >= 3 && infoParts[2] != "0" {
				// Convert bytes to MB
				if memBytes, err := strconv.ParseInt(infoParts[2], 10, 64); err == nil {
					metrics["memory_limit_mb"] = memBytes / 1024 / 1024
				}
			}
		}
	}

	helpers.WriteJSONResponse(w, http.StatusOK, metrics)
}

// GetPostgreSQLMetrics handles GET /api/v1/admin/metrics/postgresql - Get PostgreSQL database metrics
func (mc *MetricsController) GetPostgreSQLMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" && role != "ADMIN" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	metrics := map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}

	// Get database size
	var dbSize string
	mc.db.WithContext(ctx).Raw("SELECT pg_size_pretty(pg_database_size(current_database())) as size").Scan(&dbSize)
	metrics["database_size"] = dbSize

	// Get database size in bytes
	var dbSizeBytes int64
	mc.db.WithContext(ctx).Raw("SELECT pg_database_size(current_database()) as size_bytes").Scan(&dbSizeBytes)
	metrics["database_size_bytes"] = dbSizeBytes

	// Get number of connections
	var activeConnections int
	mc.db.WithContext(ctx).Raw("SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&activeConnections)
	metrics["active_connections"] = activeConnections

	var totalConnections int
	mc.db.WithContext(ctx).Raw("SELECT count(*) FROM pg_stat_activity").Scan(&totalConnections)
	metrics["total_connections"] = totalConnections

	// Get max connections
	var maxConnections int
	mc.db.WithContext(ctx).Raw("SELECT setting::int FROM pg_settings WHERE name = 'max_connections'").Scan(&maxConnections)
	metrics["max_connections"] = maxConnections

	// Get table counts
	var tableCounts []map[string]interface{}
	mc.db.WithContext(ctx).Raw(`
		SELECT 
			schemaname,
			tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
			pg_total_relation_size(schemaname||'.'||tablename) as size_bytes,
			n_live_tup as row_count
		FROM pg_tables t
		LEFT JOIN pg_stat_user_tables s ON t.tablename = s.relname
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
	`).Scan(&tableCounts)
	metrics["tables"] = tableCounts

	// Get cache hit ratio
	var cacheHitRatio float64
	mc.db.WithContext(ctx).Raw(`
		SELECT 
			round(sum(heap_blks_hit)::numeric / NULLIF(sum(heap_blks_hit) + sum(heap_blks_read), 0) * 100, 2) as ratio
		FROM pg_statio_user_tables
	`).Scan(&cacheHitRatio)
	metrics["cache_hit_ratio"] = cacheHitRatio

	// Get index usage
	var indexUsage []map[string]interface{}
	mc.db.WithContext(ctx).Raw(`
		SELECT 
			schemaname,
			tablename,
			indexname,
			pg_size_pretty(pg_relation_size(indexrelid)) as size,
			idx_scan as scans
		FROM pg_stat_user_indexes
		WHERE schemaname = 'public'
		ORDER BY pg_relation_size(indexrelid) DESC
		LIMIT 10
	`).Scan(&indexUsage)
	metrics["indexes"] = indexUsage

	helpers.WriteJSONResponse(w, http.StatusOK, metrics)
}

// GetMetrics handles GET /api/v1/admin/metrics - Get all metrics (Docker + PostgreSQL)
func (mc *MetricsController) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify admin role
	role := getUserRoleFromContext(ctx)
	if role != "admin" && role != "ADMIN" {
		helpers.WriteErrorResponse(w, http.StatusForbidden, "forbidden", "Admin access required")
		return
	}

	// Get both Docker and PostgreSQL metrics
	dockerMetrics := make(map[string]interface{})
	postgresMetrics := make(map[string]interface{})

	// Get Docker metrics
	containerName := "oms_postgres"
	cmd := exec.CommandContext(ctx, "docker", "stats", containerName, "--no-stream", "--format", "{{.CPUPerc}},{{.MemUsage}},{{.MemPerc}},{{.NetIO}},{{.BlockIO}}")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		stats := strings.TrimSpace(string(output))
		parts := strings.Split(stats, ",")
		if len(parts) >= 5 {
			dockerMetrics["cpu_percent"] = strings.TrimSpace(parts[0])
			dockerMetrics["memory_usage"] = strings.TrimSpace(parts[1])
			dockerMetrics["memory_percent"] = strings.TrimSpace(parts[2])
			dockerMetrics["network_io"] = strings.TrimSpace(parts[3])
			dockerMetrics["block_io"] = strings.TrimSpace(parts[4])
		}
	}

	// Get container status
	infoCmd := exec.CommandContext(ctx, "docker", "inspect", containerName, "--format", "{{.State.Status}},{{.State.StartedAt}}")
	infoOutput, err := infoCmd.Output()
	if err == nil && len(infoOutput) > 0 {
		infoParts := strings.Split(strings.TrimSpace(string(infoOutput)), ",")
		if len(infoParts) >= 2 {
			dockerMetrics["status"] = infoParts[0]
			dockerMetrics["started_at"] = infoParts[1]
		}
	}

	// Get PostgreSQL metrics
	var dbSize string
	mc.db.WithContext(ctx).Raw("SELECT pg_size_pretty(pg_database_size(current_database())) as size").Scan(&dbSize)
	postgresMetrics["database_size"] = dbSize

	var dbSizeBytes int64
	mc.db.WithContext(ctx).Raw("SELECT pg_database_size(current_database()) as size_bytes").Scan(&dbSizeBytes)
	postgresMetrics["database_size_bytes"] = dbSizeBytes

	var activeConnections int
	mc.db.WithContext(ctx).Raw("SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&activeConnections)
	postgresMetrics["active_connections"] = activeConnections

	var totalConnections int
	mc.db.WithContext(ctx).Raw("SELECT count(*) FROM pg_stat_activity").Scan(&totalConnections)
	postgresMetrics["total_connections"] = totalConnections

	var maxConnections int
	mc.db.WithContext(ctx).Raw("SELECT setting::int FROM pg_settings WHERE name = 'max_connections'").Scan(&maxConnections)
	postgresMetrics["max_connections"] = maxConnections

	var cacheHitRatio float64
	mc.db.WithContext(ctx).Raw(`
		SELECT 
			round(sum(heap_blks_hit)::numeric / NULLIF(sum(heap_blks_hit) + sum(heap_blks_read), 0) * 100, 2) as ratio
		FROM pg_statio_user_tables
	`).Scan(&cacheHitRatio)
	postgresMetrics["cache_hit_ratio"] = cacheHitRatio

	// Combine metrics
	response := map[string]interface{}{
		"timestamp":      time.Now().Unix(),
		"docker":         dockerMetrics,
		"postgresql":     postgresMetrics,
	}

	helpers.WriteJSONResponse(w, http.StatusOK, response)
}


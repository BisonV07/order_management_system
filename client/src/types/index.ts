// Order types
export type OrderStatus = 'ORDERED' | 'SHIPPED' | 'DELIVERED' | 'CANCELLED'

export interface Order {
  id: string
  user_id: number
  product_id: string
  quantity: number
  current_status: OrderStatus
  metadata?: Record<string, any> // Shipping address and other order metadata
  created_at: string
  updated_at: string
}

export interface Product {
  id: string
  sku: string
  name: string
  price: number
  inventory?: number // Stock quantity
  metadata: Record<string, any>
}

export interface Inventory {
  product_id: string
  quantity: number
  last_updated: string
}

// API Request/Response types
export interface CreateOrderRequest {
  product_id: string // UUID as string
  quantity: number
  shipping_address?: {
    street?: string
    city?: string
    state?: string
    zip_code?: string
    country?: string
    [key: string]: any
  }
}

export interface CreateOrderResponse {
  order_id: string
  current_status: OrderStatus
  message: string
}

export interface UpdateOrderStatusRequest {
  current_status: OrderStatus
}

export interface UpdateOrderStatusResponse {
  order_id: string
  previous_status: OrderStatus
  current_status: OrderStatus
  updated_by: number
  updated_at: string
}

export interface OrderHistory {
  order_id: string
  previous_status: OrderStatus
  new_status: OrderStatus
  updated_by: number
  updated_at: string
}

// Auth types
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user_id: number
  role: string
}

export interface SignupRequest {
  username: string
  password: string
}

export interface SignupResponse {
  message: string
  user_id: number
}

// Admin types
export interface CreateProductRequest {
  sku: string
  name: string
  price: number
  metadata?: Record<string, any>
}

export interface CreateProductResponse {
  product_id: string
  message: string
}

export interface UpdateProductRequest {
  sku?: string
  name?: string
  price?: number
  metadata?: Record<string, any>
}

export interface UpdateInventoryRequest {
  product_id: string
  quantity: number
}

export interface UpdateInventoryResponse {
  product_id: string
  quantity: number
  message: string
}

// Metrics types
export interface DockerMetrics {
  container_name?: string
  status?: string
  started_at?: string
  cpu_percent?: string
  memory_usage?: string
  memory_percent?: string
  memory_limit_mb?: number
  network_io?: string
  block_io?: string
  error?: string
  timestamp?: number
}

export interface PostgreSQLMetrics {
  database_size?: string
  database_size_bytes?: number
  active_connections?: number
  total_connections?: number
  max_connections?: number
  cache_hit_ratio?: number
  tables?: Array<{
    schemaname: string
    tablename: string
    size: string
    size_bytes: number
    row_count: number
  }>
  indexes?: Array<{
    schemaname: string
    tablename: string
    indexname: string
    size: string
    scans: number
  }>
  timestamp?: number
}

export interface SystemMetrics {
  timestamp: number
  docker: DockerMetrics
  postgresql: PostgreSQLMetrics
}


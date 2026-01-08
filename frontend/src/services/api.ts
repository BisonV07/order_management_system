import axios from 'axios'
import type {
  CreateOrderRequest,
  CreateOrderResponse,
  UpdateOrderStatusRequest,
  UpdateOrderStatusResponse,
  Product,
  Order,
  OrderHistory,
  LoginResponse,
  SignupRequest,
  SignupResponse,
  CreateProductRequest,
  CreateProductResponse,
  UpdateProductRequest,
  UpdateInventoryRequest,
  UpdateInventoryResponse,
  SystemMetrics,
  DockerMetrics,
  PostgreSQLMetrics,
} from '../types'

// Use Vite proxy in development - MUST use relative path for proxy to work
// Force relative path to ensure proxy is used (fixes CORS issues)
const envApiUrl = import.meta.env.VITE_API_BASE_URL
const API_BASE_URL = envApiUrl && !envApiUrl.startsWith('http') 
  ? envApiUrl 
  : '/api/v1' // Always use proxy path in development

// Debug: Log the API base URL to verify it's using proxy (only once on load)
if (API_BASE_URL.startsWith('http')) {
  console.warn('⚠️ WARNING: Using full URL instead of proxy! This may cause CORS issues.')
  console.warn('⚠️ Should be "/api/v1" to use Vite proxy')
}

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
})

// Add response interceptor for better error handling
apiClient.interceptors.response.use(
  (response) => {
    // Only log errors, not successful responses (to reduce console clutter)
    return response
  },
  (error) => {
    // Only log errors, with concise format
    console.error('❌', error.config?.method?.toUpperCase(), error.config?.url, '→', error.response?.status || error.code)
    if (error.response?.data?.message) {
      console.error('   Message:', error.response.data.message)
    }
    
    // Improve error messages
    if (error.code === 'ERR_NETWORK' || error.message === 'Network Error') {
      error.message = 'Cannot connect to server. Please make sure the backend is running on http://localhost:8080'
    }
    return Promise.reject(error)
  }
)

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  // Reduced logging - only log important details
  if (config.url?.includes('/auth/login') || config.url?.includes('/orders')) {
    console.log('🔵', config.method?.toUpperCase(), config.url)
  }
  
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth service functions
export const authService = {
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', {
      username,
      password,
    })
    return response.data
  },

  signup: async (username: string, password: string): Promise<SignupResponse> => {
    const response = await apiClient.post<SignupResponse>('/auth/signup', {
      username,
      password,
    })
    return response.data
  },
}

// API service functions
export const orderService = {
  createOrder: async (data: CreateOrderRequest): Promise<CreateOrderResponse> => {
    const response = await apiClient.post<CreateOrderResponse>('/orders', data)
    return response.data
  },

  updateOrderStatus: async (
    orderId: string,
    data: UpdateOrderStatusRequest
  ): Promise<UpdateOrderStatusResponse> => {
    const response = await apiClient.patch<UpdateOrderStatusResponse>(
      `/orders/${orderId}`,
      data
    )
    return response.data
  },

  getOrder: async (orderId: string): Promise<Order> => {
    const response = await apiClient.get<Order>(`/orders/${orderId}`)
    return response.data
  },

  getOrders: async (): Promise<Order[]> => {
    const response = await apiClient.get<Order[]>('/orders')
    return response.data
  },

  getOrderHistory: async (orderId: string): Promise<OrderHistory[]> => {
    const response = await apiClient.get<OrderHistory[]>(`/orders/${orderId}/history`)
    return response.data
  },
}

export const productService = {
  getProducts: async (): Promise<Product[]> => {
    const response = await apiClient.get<Product[]>('/products')
    return response.data
  },

  getProduct: async (productId: string): Promise<Product> => {
    const response = await apiClient.get<Product>(`/products/${productId}`)
    return response.data
  },
}

export const adminService = {
  createProduct: async (data: CreateProductRequest): Promise<CreateProductResponse> => {
    const response = await apiClient.post<CreateProductResponse>('/admin/products', data)
    return response.data
  },

  updateProduct: async (productId: string, data: UpdateProductRequest): Promise<Product> => {
    const response = await apiClient.put<Product>(`/admin/products/${productId}`, data)
    return response.data
  },

  deleteProduct: async (productId: string): Promise<{ message: string; product_id: string }> => {
    const response = await apiClient.delete<{ message: string; product_id: string }>(`/admin/products/${productId}`)
    return response.data
  },

  updateInventory: async (data: UpdateInventoryRequest): Promise<UpdateInventoryResponse> => {
    const response = await apiClient.put<UpdateInventoryResponse>('/admin/inventory', data)
    return response.data
  },

  getMetrics: async (): Promise<SystemMetrics> => {
    const response = await apiClient.get<SystemMetrics>('/admin/metrics')
    return response.data
  },

  getDockerMetrics: async (): Promise<DockerMetrics> => {
    const response = await apiClient.get<DockerMetrics>('/admin/metrics/docker')
    return response.data
  },

  getPostgreSQLMetrics: async (): Promise<PostgreSQLMetrics> => {
    const response = await apiClient.get<PostgreSQLMetrics>('/admin/metrics/postgresql')
    return response.data
  },
}

export default apiClient


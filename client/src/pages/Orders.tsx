import { useState, useEffect, useMemo } from 'react'
import { useAuth } from '../context/AuthContext'
import { orderService, productService } from '../services/api'
import type { CreateOrderRequest, UpdateOrderStatusRequest, OrderStatus, Order, OrderHistory, Product } from '../types'
import '../App.css'

// Trie Node for efficient prefix search
class TrieNode {
  children: Map<string, TrieNode>
  isEnd: boolean
  productIds: Set<string>

  constructor() {
    this.children = new Map()
    this.isEnd = false
    this.productIds = new Set()
  }
}

// Trie data structure for product name search
class Trie {
  root: TrieNode

  constructor() {
    this.root = new TrieNode()
  }

  insert(word: string, productId: string) {
    const normalizedWord = word.toLowerCase()
    let node = this.root
    
    for (const char of normalizedWord) {
      if (!node.children.has(char)) {
        node.children.set(char, new TrieNode())
      }
      node = node.children.get(char)!
    }
    
    node.isEnd = true
    node.productIds.add(productId)
  }

  search(prefix: string): Set<string> {
    const normalizedPrefix = prefix.toLowerCase()
    let node = this.root
    
    // Traverse to the prefix node
    for (const char of normalizedPrefix) {
      if (!node.children.has(char)) {
        return new Set() // Prefix not found
      }
      node = node.children.get(char)!
    }
    
    // Collect all product IDs from this node and its children
    const productIds = new Set<string>()
    this.collectAllProductIds(node, productIds)
    return productIds
  }

  private collectAllProductIds(node: TrieNode, productIds: Set<string>) {
    if (node.isEnd) {
      node.productIds.forEach(id => productIds.add(id))
    }
    
    for (const child of node.children.values()) {
      this.collectAllProductIds(child, productIds)
    }
  }
}

const Orders = () => {
  const { role } = useAuth()
  const [orderId, setOrderId] = useState('')
  const [newStatus, setNewStatus] = useState<OrderStatus>('SHIPPED')
  const [message, setMessage] = useState<string | null>(null)

  const [productId, setProductId] = useState('550e8400-e29b-41d4-a716-446655440000') // Sample product ID
  const [quantity, setQuantity] = useState(0)
  
  // Shipping address state
  const [shippingAddress, setShippingAddress] = useState({
    street: '',
    city: '',
    state: '',
    zip_code: '',
    country: ''
  })

  // Order tracking state
  const [orders, setOrders] = useState<Order[]>([])
  const [loadingOrders, setLoadingOrders] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [orderHistory, setOrderHistory] = useState<OrderHistory[]>([])
  const [loadingHistory, setLoadingHistory] = useState(false)

  // Search state
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<OrderStatus | 'ALL'>('ALL')
  const [products, setProducts] = useState<Product[]>([])
  const [productMap, setProductMap] = useState<Map<string, Product>>(new Map())

  // Build Trie for product name search
  const productTrie = useMemo(() => {
    const trie = new Trie()
    products.forEach(product => {
      // Insert product name into trie
      trie.insert(product.name, product.id)
      // Also insert SKU for search
      trie.insert(product.sku, product.id)
    })
    return trie
  }, [products])

  // Load orders on mount
  useEffect(() => {
    loadOrders()
    // Load products for search functionality
    loadProducts()
  }, [])

  // Load products for search
  const loadProducts = async () => {
    try {
      const data = await productService.getProducts()
      setProducts(data)
      // Create a map for quick product lookup
      const map = new Map<string, Product>()
      data.forEach(product => {
        map.set(product.id, product)
      })
      setProductMap(map)
    } catch (err: any) {
      console.error('Failed to load products for search:', err)
    }
  }

  // Update newStatus when order is selected to show valid transitions
  useEffect(() => {
    if (selectedOrder) {
      // Set default next status based on current status and user role
      if (selectedOrder.current_status === 'ORDERED') {
        // Regular users default to CANCELLED, admin defaults to SHIPPED
        if (role === 'admin' || role === 'ADMIN') {
          setNewStatus('SHIPPED')
        } else {
          setNewStatus('CANCELLED')
        }
      } else if (selectedOrder.current_status === 'SHIPPED') {
        setNewStatus('DELIVERED')
      }
      // For DELIVERED or CANCELLED, keep current selection
    } else if (orderId) {
      // If orderId is manually entered, try to find the order
      const order = orders.find(o => o.id === orderId)
      if (order) {
        if (order.current_status === 'ORDERED') {
          // Regular users default to CANCELLED, admin defaults to SHIPPED
          if (role === 'admin' || role === 'ADMIN') {
            setNewStatus('SHIPPED')
          } else {
            setNewStatus('CANCELLED')
          }
        } else if (order.current_status === 'SHIPPED') {
          setNewStatus('DELIVERED')
        }
      } else {
        // Order not found in list, default based on role
        if (role === 'admin' || role === 'ADMIN') {
          setNewStatus('SHIPPED')
        } else {
          setNewStatus('CANCELLED')
        }
      }
    }
  }, [selectedOrder, orderId, orders, role])

  const loadOrders = async () => {
    try {
      setLoadingOrders(true)
      const data = await orderService.getOrders()
      setOrders(data)
    } catch (err: any) {
      console.error('Failed to load orders:', err)
      setMessage(`Error loading orders: ${err?.response?.data?.message || err?.message}`)
    } finally {
      setLoadingOrders(false)
    }
  }

  const loadOrderHistory = async (orderId: string) => {
    try {
      setLoadingHistory(true)
      const data = await orderService.getOrderHistory(orderId)
      setOrderHistory(data)
    } catch (err: any) {
      console.error('Failed to load order history:', err)
      setMessage(`Error loading order history: ${err?.response?.data?.message || err?.message}`)
    } finally {
      setLoadingHistory(false)
    }
  }

  const handleCreateOrder = async () => {
    // Validate quantity
    if (quantity <= 0 || quantity > 100000000) {
      setMessage('❌ Error: Quantity must be between 1 and 100,000,000')
      return
    }
    
    // Validate shipping address (at least city and zip_code required)
    if (!shippingAddress.city || !shippingAddress.zip_code) {
      setMessage('❌ Error: City and Zip Code are required for shipping address')
      return
    }
    
    try {
      // Build shipping address object (only include non-empty fields)
      const address: Record<string, string> = {}
      if (shippingAddress.street) address.street = shippingAddress.street
      if (shippingAddress.city) address.city = shippingAddress.city
      if (shippingAddress.state) address.state = shippingAddress.state
      if (shippingAddress.zip_code) address.zip_code = shippingAddress.zip_code
      if (shippingAddress.country) address.country = shippingAddress.country
      
      const request: CreateOrderRequest = {
        product_id: productId,
        quantity: quantity,
        shipping_address: address
      }
      const response = await orderService.createOrder(request)
      setMessage(`✅ Order created successfully! Order ID: ${response.order_id}`)
      setProductId('')
      setQuantity(0)
      setShippingAddress({
        street: '',
        city: '',
        state: '',
        zip_code: '',
        country: ''
      })
      // Reload orders list
      loadOrders()
      // Trigger product refresh event for Dashboard
      window.dispatchEvent(new Event('refresh-products'))
    } catch (err: any) {
      const errorMsg = err?.response?.data?.message || err?.message || 'Failed to create order'
      setMessage(`❌ Error: ${errorMsg}`)
      console.error(err)
    }
  }

  const handleUpdateStatus = async () => {
    if (!orderId) {
      setMessage('Please enter an order ID')
      return
    }

    try {
      const request: UpdateOrderStatusRequest = {
        current_status: newStatus,
      }
      const response = await orderService.updateOrderStatus(orderId, request)
      setMessage(`✅ Order ${response.order_id} updated from ${response.previous_status} to ${response.current_status}`)
      setOrderId('')
      // Reload orders list and history if viewing this order
      loadOrders()
      if (selectedOrder?.id === orderId) {
        loadOrderHistory(orderId)
      }
      // If order was cancelled, refresh products to restore inventory
      if (newStatus === 'CANCELLED') {
        window.dispatchEvent(new Event('refresh-products'))
      }
    } catch (err: any) {
      const errorMsg = err?.response?.data?.message || err?.message || 'Failed to update order status'
      setMessage(`❌ Error: ${errorMsg}`)
      console.error(err)
    }
  }

  const handleViewOrder = async (order: Order) => {
    setSelectedOrder(order)
    setOrderId(order.id) // Pre-fill order ID for status update
    await loadOrderHistory(order.id)
  }

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case 'ORDERED': return 'var(--info)'
      case 'SHIPPED': return 'var(--warning)'
      case 'DELIVERED': return 'var(--success)'
      case 'CANCELLED': return 'var(--danger)'
      default: return 'var(--gray)'
    }
  }

  const getStatusBadge = (status: OrderStatus) => {
    const colors: Record<OrderStatus, string> = {
      'ORDERED': 'badge-info',
      'SHIPPED': 'badge-warning',
      'DELIVERED': 'badge-success',
      'CANCELLED': 'badge-danger'
    }
    return colors[status] || 'badge-info'
  }

  // Filter orders based on search query and status filter
  const filteredOrders = useMemo(() => {
    let filtered = orders

    // Apply status filter
    if (statusFilter !== 'ALL') {
      filtered = filtered.filter(order => order.current_status === statusFilter)
    }

    // Apply search query filter
    if (searchQuery.trim()) {
      const query = searchQuery.trim().toLowerCase()
      const searchResults: Order[] = []

      // Search by order ID (partial match)
      const orderIdMatches = filtered.filter(order => 
        order.id.toLowerCase().includes(query)
      )

      // Search by product name using Trie (prefix match)
      const productIdsFromTrie = productTrie.search(query)
      const productNameMatches = filtered.filter(order =>
        productIdsFromTrie.has(order.product_id)
      )

      // Combine and deduplicate results
      const allMatches = [...orderIdMatches, ...productNameMatches]
      const seen = new Set<string>()
      
      allMatches.forEach(order => {
        if (!seen.has(order.id)) {
          seen.add(order.id)
          searchResults.push(order)
        }
      })

      filtered = searchResults
    }

    return filtered
  }, [orders, searchQuery, statusFilter, productTrie])

  return (
    <div className="app">
      <div className="container" style={{ paddingBottom: '200px' }}>
        <div className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h1 style={{ 
            margin: 0,
            background: 'linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent'
          }}>
            📦 Order Management
          </h1>
          <p style={{ margin: '8px 0 0', color: 'var(--gray)', fontSize: '14px' }}>
            {role === 'admin' || role === 'ADMIN' ? 'View and manage all orders' : 'Track and manage your orders'}
          </p>
        </div>

        {/* Create Order Section - Hidden for admin */}
        {role !== 'admin' && role !== 'ADMIN' && (
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px' }}>
            🛒 Place New Order
          </h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', maxWidth: '400px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Product ID (UUID):
              </label>
              <input
                type="text"
                value={productId}
                onChange={(e) => setProductId(e.target.value)}
                placeholder="Enter product UUID"
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Quantity (0 - 100,000,000):
              </label>
              <input
                type="number"
                value={quantity}
                onChange={(e) => {
                  const val = parseInt(e.target.value) || 0
                  if (val >= 0 && val <= 100000000) {
                    setQuantity(val)
                  }
                }}
                min="0"
                max="100000000"
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box',
                  MozAppearance: 'textfield', // Remove spinner in Firefox
                  WebkitAppearance: 'none', // Remove spinner in Chrome/Safari
                  appearance: 'none' // Remove spinner in modern browsers
                }}
                onWheel={(e) => e.currentTarget.blur()} // Prevent scroll from changing value
                onKeyDown={(e) => {
                  // Prevent arrow keys from changing value
                  if (e.key === 'ArrowUp' || e.key === 'ArrowDown') {
                    e.preventDefault()
                  }
                }}
              />
            </div>
            
            {/* Shipping Address Section */}
            <div style={{ marginTop: '20px', paddingTop: '20px', borderTop: '1px solid var(--gray-lighter)' }}>
              <h3 style={{ marginTop: 0, marginBottom: '15px', fontSize: '1.1rem', fontWeight: '600', color: 'var(--dark)' }}>
                📍 Shipping Address
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500', fontSize: '14px' }}>
                    Street Address:
                  </label>
                  <input
                    type="text"
                    value={shippingAddress.street}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, street: e.target.value })}
                    placeholder="123 Main St"
                    style={{
                      width: '100%',
                      padding: '8px',
                      border: '1px solid #ddd',
                      borderRadius: '4px',
                      boxSizing: 'border-box'
                    }}
                  />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '12px' }}>
                  <div>
                    <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500', fontSize: '14px' }}>
                      City <span style={{ color: 'var(--danger)' }}>*</span>:
                    </label>
                    <input
                      type="text"
                      value={shippingAddress.city}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, city: e.target.value })}
                      placeholder="New York"
                      required
                      style={{
                        width: '100%',
                        padding: '8px',
                        border: '1px solid #ddd',
                        borderRadius: '4px',
                        boxSizing: 'border-box'
                      }}
                    />
                  </div>
                  <div>
                    <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500', fontSize: '14px' }}>
                      Zip Code <span style={{ color: 'var(--danger)' }}>*</span>:
                    </label>
                    <input
                      type="text"
                      value={shippingAddress.zip_code}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, zip_code: e.target.value })}
                      placeholder="10001"
                      required
                      style={{
                        width: '100%',
                        padding: '8px',
                        border: '1px solid #ddd',
                        borderRadius: '4px',
                        boxSizing: 'border-box'
                      }}
                    />
                  </div>
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                  <div>
                    <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500', fontSize: '14px' }}>
                      State:
                    </label>
                    <input
                      type="text"
                      value={shippingAddress.state}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, state: e.target.value })}
                      placeholder="NY"
                      style={{
                        width: '100%',
                        padding: '8px',
                        border: '1px solid #ddd',
                        borderRadius: '4px',
                        boxSizing: 'border-box'
                      }}
                    />
                  </div>
                  <div>
                    <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500', fontSize: '14px' }}>
                      Country:
                    </label>
                    <input
                      type="text"
                      value={shippingAddress.country}
                      onChange={(e) => setShippingAddress({ ...shippingAddress, country: e.target.value })}
                      placeholder="USA"
                      style={{
                        width: '100%',
                        padding: '8px',
                        border: '1px solid #ddd',
                        borderRadius: '4px',
                        boxSizing: 'border-box'
                      }}
                    />
                  </div>
                </div>
                <p style={{ margin: 0, fontSize: '12px', color: 'var(--gray)', fontStyle: 'italic' }}>
                  * Required fields. All orders ship from our single warehouse.
                </p>
              </div>
            </div>
            
            <button
              onClick={handleCreateOrder}
              className="btn btn-primary"
              style={{ fontSize: '16px', padding: '12px 24px', marginTop: '10px' }}
            >
              ✨ Place Order
            </button>
          </div>
        </section>
        )}

        {/* Track Existing Orders Section */}
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <div className="flex-between" style={{ marginBottom: '20px' }}>
            <h2 style={{ margin: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)' }}>
              {role === 'admin' || role === 'ADMIN' ? '📋 All Orders' : '📋 My Orders'}
            </h2>
            <button
              onClick={loadOrders}
              className="btn btn-primary"
              style={{ fontSize: '14px', padding: '10px 20px' }}
            >
              🔄 Refresh
            </button>
          </div>
          
          {/* Search Bar and Status Filter */}
          {orders.length > 0 && (
            <div style={{ marginBottom: '20px' }}>
              <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'flex-start' }}>
                {/* Search Input */}
                <div style={{ position: 'relative', flex: '1', minWidth: '250px', maxWidth: '500px' }}>
                  <input
                    type="text"
                    placeholder="🔍 Search by order ID or product name..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    style={{
                      width: '100%',
                      padding: '12px 16px 12px 44px',
                      border: '1px solid var(--gray-lighter)',
                      borderRadius: 'var(--radius)',
                      fontSize: '14px',
                      boxSizing: 'border-box',
                      background: 'rgba(255, 255, 255, 0.9)',
                      transition: 'var(--transition)'
                    }}
                    onFocus={(e) => {
                      e.currentTarget.style.borderColor = 'var(--primary)'
                      e.currentTarget.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.1)'
                    }}
                    onBlur={(e) => {
                      e.currentTarget.style.borderColor = 'var(--gray-lighter)'
                      e.currentTarget.style.boxShadow = 'none'
                    }}
                  />
                  <span style={{
                    position: 'absolute',
                    left: '16px',
                    top: '50%',
                    transform: 'translateY(-50%)',
                    fontSize: '18px',
                    pointerEvents: 'none'
                  }}>
                    🔍
                  </span>
                  {searchQuery && (
                    <button
                      onClick={() => setSearchQuery('')}
                      style={{
                        position: 'absolute',
                        right: '8px',
                        top: '50%',
                        transform: 'translateY(-50%)',
                        background: 'transparent',
                        border: 'none',
                        cursor: 'pointer',
                        fontSize: '18px',
                        padding: '4px 8px',
                        color: 'var(--gray)',
                        borderRadius: 'var(--radius)',
                        transition: 'var(--transition)'
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'
                        e.currentTarget.style.color = 'var(--danger)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.background = 'transparent'
                        e.currentTarget.style.color = 'var(--gray)'
                      }}
                      title="Clear search"
                    >
                      ✕
                    </button>
                  )}
                </div>

                {/* Status Filter Dropdown */}
                <div style={{ minWidth: '180px' }}>
                  <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value as OrderStatus | 'ALL')}
                    style={{
                      width: '100%',
                      padding: '12px 16px 12px 40px',
                      border: '1px solid var(--gray-lighter)',
                      borderRadius: 'var(--radius)',
                      fontSize: '14px',
                      boxSizing: 'border-box',
                      background: 'rgba(255, 255, 255, 0.9)',
                      cursor: 'pointer',
                      appearance: 'none',
                      backgroundImage: 'url("data:image/svg+xml,%3Csvg xmlns=\'http://www.w3.org/2000/svg\' width=\'12\' height=\'12\' viewBox=\'0 0 12 12\'%3E%3Cpath fill=\'%23666\' d=\'M6 9L1 4h10z\'/%3E%3C/svg%3E")',
                      backgroundRepeat: 'no-repeat',
                      backgroundPosition: 'right 12px center',
                      paddingRight: '36px',
                      transition: 'var(--transition)'
                    }}
                    onFocus={(e) => {
                      e.currentTarget.style.borderColor = 'var(--primary)'
                      e.currentTarget.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.1)'
                    }}
                    onBlur={(e) => {
                      e.currentTarget.style.borderColor = 'var(--gray-lighter)'
                      e.currentTarget.style.boxShadow = 'none'
                    }}
                  >
                    <option value="ALL">📋 All Statuses</option>
                    <option value="ORDERED">📦 Ordered</option>
                    <option value="SHIPPED">🚚 Shipped</option>
                    <option value="DELIVERED">✅ Delivered</option>
                    <option value="CANCELLED">❌ Cancelled</option>
                  </select>
                </div>
              </div>
              
              {/* Filter Results Summary */}
              {(searchQuery || statusFilter !== 'ALL') && (
                <div style={{ 
                  marginTop: '12px', 
                  display: 'flex', 
                  gap: '12px',
                  flexWrap: 'wrap',
                  alignItems: 'center'
                }}>
                  <p style={{ 
                    margin: 0,
                    fontSize: '13px', 
                    color: 'var(--gray)',
                    fontWeight: '500'
                  }}>
                    Found <strong style={{ color: 'var(--primary)' }}>{filteredOrders.length}</strong> {filteredOrders.length === 1 ? 'order' : 'orders'}
                    {searchQuery && ` matching "${searchQuery}"`}
                    {statusFilter !== 'ALL' && ` with status "${statusFilter}"`}
                  </p>
                  {(searchQuery || statusFilter !== 'ALL') && (
                    <button
                      onClick={() => {
                        setSearchQuery('')
                        setStatusFilter('ALL')
                      }}
                      style={{
                        padding: '6px 12px',
                        fontSize: '12px',
                        background: 'rgba(239, 68, 68, 0.1)',
                        color: 'var(--danger)',
                        border: '1px solid rgba(239, 68, 68, 0.2)',
                        borderRadius: 'var(--radius)',
                        cursor: 'pointer',
                        fontWeight: '500',
                        transition: 'var(--transition)'
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.background = 'rgba(239, 68, 68, 0.2)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)'
                      }}
                    >
                      Clear Filters
                    </button>
                  )}
                </div>
              )}
            </div>
          )}
          
          {loadingOrders ? (
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <div className="loading" style={{ margin: '0 auto 16px' }}></div>
              <p style={{ color: 'var(--gray)' }}>Loading orders...</p>
            </div>
          ) : orders.length === 0 ? (
            <div className="card" style={{ textAlign: 'center', padding: '60px 20px', background: 'rgba(99, 102, 241, 0.05)' }}>
              <div style={{ fontSize: '48px', marginBottom: '16px' }}>📭</div>
              <h3 style={{ color: 'var(--gray)', marginBottom: '8px' }}>No orders found</h3>
              <p style={{ color: 'var(--gray-light)', fontSize: '14px' }}>
                {role === 'admin' || role === 'ADMIN' ? 'Orders will appear here once users create them.' : 'Create your first order above!'}
              </p>
            </div>
          ) : filteredOrders.length === 0 && (searchQuery || statusFilter !== 'ALL') ? (
            <div className="card" style={{ textAlign: 'center', padding: '60px 20px', background: 'rgba(99, 102, 241, 0.05)' }}>
              <div style={{ fontSize: '48px', marginBottom: '16px' }}>🔍</div>
              <h3 style={{ color: 'var(--gray)', marginBottom: '8px' }}>No matching orders</h3>
              <p style={{ color: 'var(--gray-light)', fontSize: '14px', marginBottom: '12px' }}>
                {searchQuery && statusFilter !== 'ALL' 
                  ? `No orders found matching "${searchQuery}" with status "${statusFilter}".`
                  : searchQuery
                  ? `No orders found matching "${searchQuery}".`
                  : `No orders found with status "${statusFilter}".`
                }
              </p>
              <button
                onClick={() => {
                  setSearchQuery('')
                  setStatusFilter('ALL')
                }}
                className="btn btn-primary"
                style={{ fontSize: '14px', padding: '8px 16px' }}
              >
                Clear Filters
              </button>
            </div>
          ) : (
            <div className="grid" style={{ gridTemplateColumns: '1fr', gap: '16px' }}>
              {filteredOrders.map((order) => {
                const product = productMap.get(order.product_id)
                return (
                  <div
                    key={order.id}
                    onClick={() => handleViewOrder(order)}
                    className="card"
                    style={{
                      cursor: 'pointer',
                      background: selectedOrder?.id === order.id 
                        ? 'linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(99, 102, 241, 0.05) 100%)'
                        : 'rgba(255, 255, 255, 0.95)',
                      border: selectedOrder?.id === order.id 
                        ? '2px solid var(--primary)' 
                        : '1px solid var(--gray-lighter)',
                      transition: 'var(--transition)'
                    }}
                    onMouseEnter={(e) => {
                      if (selectedOrder?.id !== order.id) {
                        e.currentTarget.style.transform = 'translateY(-2px)'
                        e.currentTarget.style.boxShadow = 'var(--shadow-lg)'
                      }
                    }}
                    onMouseLeave={(e) => {
                      if (selectedOrder?.id !== order.id) {
                        e.currentTarget.style.transform = 'translateY(0)'
                        e.currentTarget.style.boxShadow = 'var(--shadow-md)'
                      }
                    }}
                  >
                    <div className="flex-between">
                      <div style={{ flex: 1 }}>
                        <div style={{ fontWeight: '700', marginBottom: '8px', fontSize: '16px', color: 'var(--dark)' }}>
                          Order: <code style={{ 
                            fontSize: '0.9em', 
                            background: 'rgba(99, 102, 241, 0.1)', 
                            padding: '4px 8px', 
                            borderRadius: '4px',
                            color: 'var(--primary)',
                            fontWeight: '600'
                          }}>{order.id.slice(0, 8)}...</code>
                        </div>
                        {role === 'admin' || role === 'ADMIN' ? (
                          <div style={{ fontSize: '0.9em', color: 'var(--gray)', marginBottom: '6px' }}>
                            👤 User ID: <strong>{order.user_id}</strong>
                          </div>
                        ) : null}
                        <div style={{ fontSize: '0.9em', color: 'var(--gray)', marginBottom: '6px' }}>
                          📦 Product: {product ? (
                            <span style={{ fontWeight: '600', color: 'var(--primary)' }}>{product.name}</span>
                          ) : (
                            <code style={{ fontSize: '0.85em', background: 'var(--gray-lighter)', padding: '2px 6px', borderRadius: '3px' }}>{order.product_id.slice(0, 8)}...</code>
                          )} | 
                          Quantity: <strong>{order.quantity}</strong>
                        </div>
                        {order.metadata && order.metadata.shipping_address && (
                          <div style={{ fontSize: '0.85em', color: 'var(--gray)', marginBottom: '6px', marginTop: '4px' }}>
                            📍 Shipping: {(() => {
                              const addr = order.metadata.shipping_address as any
                              const parts = []
                              if (addr.street) parts.push(addr.street)
                              if (addr.city) parts.push(addr.city)
                              if (addr.state) parts.push(addr.state)
                              if (addr.zip_code) parts.push(addr.zip_code)
                              if (addr.country) parts.push(addr.country)
                              return parts.length > 0 ? parts.join(', ') : 'Address not provided'
                            })()}
                          </div>
                        )}
                        <div style={{ fontSize: '0.85em', color: 'var(--gray-light)' }}>
                          🕒 {new Date(order.created_at).toLocaleString()}
                        </div>
                      </div>
                      <span className={`badge ${getStatusBadge(order.current_status)}`} style={{ fontSize: '12px', padding: '8px 16px' }}>
                        {order.current_status}
                      </span>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </section>

        {/* Order History Section */}
        {selectedOrder && (
          <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
            <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px' }}>
              📜 Order History & Tracking
            </h2>
            <div className="card" style={{ 
              marginBottom: '20px', 
              background: 'rgba(99, 102, 241, 0.05)',
              border: '1px solid rgba(99, 102, 241, 0.2)'
            }}>
              <div style={{ marginBottom: '8px' }}>
                <strong style={{ color: 'var(--dark)' }}>Order ID:</strong>{' '}
                <code style={{ fontSize: '0.9em', background: 'rgba(99, 102, 241, 0.1)', padding: '4px 8px', borderRadius: '4px', color: 'var(--primary)' }}>
                  {selectedOrder.id}
                </code>
              </div>
              <div style={{ marginBottom: '8px' }}>
                <strong style={{ color: 'var(--dark)' }}>Current Status:</strong>{' '}
                <span className={`badge ${getStatusBadge(selectedOrder.current_status)}`} style={{ fontSize: '12px', padding: '4px 12px' }}>
                  {selectedOrder.current_status}
                </span>
              </div>
              {selectedOrder.metadata && selectedOrder.metadata.shipping_address && (
                <div style={{ marginTop: '12px', paddingTop: '12px', borderTop: '1px solid rgba(99, 102, 241, 0.2)' }}>
                  <strong style={{ color: 'var(--dark)', display: 'block', marginBottom: '6px' }}>📍 Shipping Address:</strong>
                  <div style={{ fontSize: '0.9em', color: 'var(--gray)', lineHeight: '1.6' }}>
                    {(() => {
                      const addr = selectedOrder.metadata.shipping_address as any
                      const parts = []
                      if (addr.street) parts.push(addr.street)
                      if (addr.city || addr.state || addr.zip_code) {
                        const cityParts = []
                        if (addr.city) cityParts.push(addr.city)
                        if (addr.state) cityParts.push(addr.state)
                        if (addr.zip_code) cityParts.push(addr.zip_code)
                        parts.push(cityParts.join(', '))
                      }
                      if (addr.country) parts.push(addr.country)
                      return parts.length > 0 ? parts.join('\n') : 'Address not provided'
                    })()}
                  </div>
                </div>
              )}
            </div>
            
            {loadingHistory ? (
              <div style={{ textAlign: 'center', padding: '40px' }}>
                <div className="loading" style={{ margin: '0 auto 16px' }}></div>
                <p style={{ color: 'var(--gray)' }}>Loading history...</p>
              </div>
            ) : orderHistory.length === 0 ? (
              <div className="card" style={{ textAlign: 'center', padding: '40px', background: 'rgba(99, 102, 241, 0.05)' }}>
                <div style={{ fontSize: '36px', marginBottom: '12px' }}>📝</div>
                <p style={{ color: 'var(--gray)', margin: 0 }}>No history available. This order was just created.</p>
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                {orderHistory.map((history, idx) => (
                  <div
                    key={idx}
                    className="card"
                    style={{
                      borderLeft: `4px solid ${getStatusColor(history.new_status as OrderStatus)}`,
                      background: 'rgba(255, 255, 255, 0.8)',
                      padding: '16px'
                    }}
                  >
                    <div style={{ fontWeight: '700', marginBottom: '8px', fontSize: '16px', color: 'var(--dark)' }}>
                      <span className={`badge ${getStatusBadge(history.previous_status as OrderStatus)}`} style={{ fontSize: '11px', padding: '4px 8px', marginRight: '8px' }}>
                        {history.previous_status}
                      </span>
                      <span style={{ fontSize: '18px', margin: '0 8px' }}>→</span>
                      <span className={`badge ${getStatusBadge(history.new_status as OrderStatus)}`} style={{ fontSize: '11px', padding: '4px 8px' }}>
                        {history.new_status}
                      </span>
                    </div>
                    <div style={{ fontSize: '0.85em', color: 'var(--gray)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <span>👤</span>
                      <span>Updated by User #{history.updated_by}</span>
                      <span style={{ margin: '0 4px' }}>•</span>
                      <span>🕒</span>
                      <span>{new Date(history.updated_at).toLocaleString()}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        )}

        {/* Update Order Status Section */}
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px' }}>
            {role === 'admin' || role === 'ADMIN' ? '⚙️ Update Order Status' : '✏️ Change Order Status'}
            {role === 'admin' || role === 'ADMIN' ? (
              <span className="badge badge-danger" style={{ marginLeft: '12px', fontSize: '12px' }}>
                (SHIPPED/DELIVERED only)
              </span>
            ) : null}
          </h2>
          {role === 'admin' && (
            <p style={{ fontSize: '0.9em', color: '#666', marginBottom: '15px' }}>
              As admin, you can update any order's status to SHIPPED or DELIVERED. Regular users can only cancel ORDERED orders (cannot update to SHIPPED).
            </p>
          )}
          {role !== 'admin' && (
            <p style={{ fontSize: '0.9em', color: '#666', marginBottom: '15px' }}>
              You can only cancel ORDERED orders. Only admin can update orders to SHIPPED or DELIVERED.
            </p>
          )}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', maxWidth: '400px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Order ID (UUID):
              </label>
              <input
                type="text"
                placeholder="Enter order UUID or click an order above"
                value={orderId}
                onChange={(e) => setOrderId(e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                New Status (FSM Validated):
              </label>
              <select
                value={newStatus}
                onChange={(e) => setNewStatus(e.target.value as OrderStatus)}
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
              >
                {(() => {
                  // Get current order status to show only valid transitions
                  const currentOrder = selectedOrder || orders.find(o => o.id === orderId)
                  const currentStatus = currentOrder?.current_status || 'ORDERED'
                  
                  // FSM Rules:
                  // ORDERED → SHIPPED or CANCELLED
                  // SHIPPED → DELIVERED
                  // DELIVERED → (no transitions)
                  // CANCELLED → (no transitions)
                  
                  const options = []
                  
                  // Admin can only update to SHIPPED or DELIVERED
                  // Regular users can cancel ORDERED orders, but SHIPPED orders cannot be cancelled
                  if (role === 'admin' || role === 'ADMIN') {
                    // Admin restrictions: only SHIPPED or DELIVERED
                    if (currentStatus === 'ORDERED') {
                      options.push(<option key="SHIPPED" value="SHIPPED">Shipped (ORDERED → SHIPPED)</option>)
                    } else if (currentStatus === 'SHIPPED') {
                      options.push(<option key="DELIVERED" value="DELIVERED">Delivered (SHIPPED → DELIVERED)</option>)
                    } else {
                      options.push(<option key="FINAL" value={currentStatus} disabled>{currentStatus} (Final state - no transitions)</option>)
                    }
                  } else {
                    // Regular user: can ONLY cancel ORDERED orders (cannot update to SHIPPED or DELIVERED)
                    if (currentStatus === 'ORDERED') {
                      options.push(<option key="CANCELLED" value="CANCELLED">Cancelled (ORDERED → CANCELLED)</option>)
                    } else if (currentStatus === 'SHIPPED') {
                      // SHIPPED orders cannot be changed by regular users - only admin can update to DELIVERED
                      options.push(<option key="DELIVERED" value="DELIVERED" disabled>Delivered (SHIPPED → DELIVERED) - Admin only</option>)
                    } else if (currentStatus === 'DELIVERED') {
                      options.push(<option key="DELIVERED" value="DELIVERED" disabled>Delivered (Final state - no transitions)</option>)
                    } else if (currentStatus === 'CANCELLED') {
                      options.push(<option key="CANCELLED" value="CANCELLED" disabled>Cancelled (Final state - no transitions)</option>)
                    }
                  }
                  
                  return options.length > 0 ? options : (
                    <>
                      <option value="SHIPPED">Shipped (ORDERED → SHIPPED)</option>
                      <option value="DELIVERED">Delivered (SHIPPED → DELIVERED)</option>
                      <option value="CANCELLED">Cancelled (Any → CANCELLED)</option>
                    </>
                  )
                })()}
              </select>
              <small style={{ color: '#666', fontSize: '0.85em', display: 'block', marginTop: '5px' }}>
                {(() => {
                  const currentOrder = selectedOrder || orders.find(o => o.id === orderId)
                  const currentStatus = currentOrder?.current_status || 'ORDERED'
                  
                  if (role === 'admin' || role === 'ADMIN') {
                    if (currentStatus === 'ORDERED') {
                      return 'Admin: ORDERED → Can update to: SHIPPED only'
                    } else if (currentStatus === 'SHIPPED') {
                      return 'Admin: SHIPPED → Can update to: DELIVERED only'
                    } else {
                      return `Admin: ${currentStatus} → Final state (no transitions allowed)`
                    }
                  } else {
                    if (currentStatus === 'ORDERED') {
                      return 'Current: ORDERED → Can only cancel (CANCELLED). Only admin can update to SHIPPED'
                    } else if (currentStatus === 'SHIPPED') {
                      return 'Current: SHIPPED → Cannot be changed. Only admin can update to DELIVERED'
                    } else if (currentStatus === 'DELIVERED' || currentStatus === 'CANCELLED') {
                      return `Current: ${currentStatus} → Final state (no transitions allowed)`
                    }
                  }
                  return 'FSM validates transitions. Invalid transitions will be rejected.'
                })()}
              </small>
            </div>
            <button
              onClick={handleUpdateStatus}
              className="btn btn-success"
              style={{ fontSize: '16px', padding: '12px 24px' }}
            >
              ✅ Update Status
            </button>
          </div>
        </section>

        {message && (
          <div className={`alert ${message.includes('❌') || message.includes('Error') ? 'alert-error' : 'alert-success'}`} style={{ marginTop: '20px' }}>
            <span>{message.includes('❌') || message.includes('Error') ? '⚠️' : '✅'}</span>
            <span>{message}</span>
          </div>
        )}
      </div>
    </div>
  )
}

export default Orders

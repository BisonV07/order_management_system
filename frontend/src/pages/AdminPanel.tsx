import { useState, useEffect } from 'react'
import { useAuth } from '../context/AuthContext'
import { productService, adminService } from '../services/api'
import type { Product, CreateProductRequest, UpdateInventoryRequest } from '../types'
import '../App.css'

const AdminPanel = () => {
  const { role } = useAuth()
  const isAdmin = role === 'ADMIN' || role === 'admin'

  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState<string | null>(null)

  // Product creation form
  const [newProduct, setNewProduct] = useState<CreateProductRequest>({
    sku: '',
    name: '',
    price: 0,
    metadata: {},
  })

  // Inventory update form
  const [inventoryUpdate, setInventoryUpdate] = useState<UpdateInventoryRequest>({
    product_id: '',
    quantity: 0,
  })

  useEffect(() => {
    if (isAdmin) {
      loadProducts()
    }
  }, [isAdmin])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const data = await productService.getProducts()
      setProducts(data)
    } catch (err: any) {
      setMessage(`Error loading products: ${err?.response?.data?.message || err?.message}`)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateProduct = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage(null)

    if (!newProduct.sku || !newProduct.name || newProduct.price <= 0) {
      setMessage('Please fill in all required fields (SKU, Name, Price > 0)')
      return
    }

    try {
      await adminService.createProduct(newProduct)
      setMessage('✅ Product created successfully!')
      setNewProduct({ sku: '', name: '', price: 0, metadata: {} })
      loadProducts()
      // Refresh products on Dashboard
      window.dispatchEvent(new CustomEvent('refresh-products'))
    } catch (err: any) {
      setMessage(`❌ Error: ${err?.response?.data?.message || err?.message || 'Failed to create product'}`)
    }
  }

  const handleUpdateInventory = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage(null)

    if (!inventoryUpdate.product_id || inventoryUpdate.quantity < 0) {
      setMessage('Please select a product and enter a valid quantity (>= 0)')
      return
    }

    try {
      await adminService.updateInventory(inventoryUpdate)
      setMessage(`✅ Inventory updated successfully! Product ${inventoryUpdate.product_id} now has ${inventoryUpdate.quantity} units`)
      setInventoryUpdate({ product_id: '', quantity: 0 })
      loadProducts()
      // Refresh products on Dashboard
      window.dispatchEvent(new CustomEvent('refresh-products'))
    } catch (err: any) {
      setMessage(`❌ Error: ${err?.response?.data?.message || err?.message || 'Failed to update inventory'}`)
    }
  }

  const handleDeleteProduct = async (productId: string, productName: string) => {
    if (!confirm(`Are you sure you want to delete "${productName}"? This action cannot be undone.`)) {
      return
    }

    setMessage(null)
    try {
      await adminService.deleteProduct(productId)
      setMessage(`✅ Product "${productName}" deleted successfully!`)
      loadProducts()
      // Refresh products on Dashboard
      window.dispatchEvent(new CustomEvent('refresh-products'))
    } catch (err: any) {
      setMessage(`❌ Error: ${err?.response?.data?.message || err?.message || 'Failed to delete product'}`)
    }
  }

  if (!isAdmin) {
    return (
      <div className="app">
        <div className="container">
          <div className="card" style={{ textAlign: 'center', padding: '60px 20px' }}>
            <div style={{ fontSize: '64px', marginBottom: '20px' }}>🚫</div>
            <h1 style={{ color: 'var(--danger)', marginBottom: '16px' }}>Access Denied</h1>
            <p style={{ color: 'var(--gray)', marginBottom: '8px' }}>You must be an admin to access this page.</p>
            <p className="badge badge-info" style={{ fontSize: '12px', marginTop: '16px', display: 'inline-block' }}>
              Your current role: {role || 'not set'}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="app">
      <div className="container" style={{ paddingBottom: '200px' }}>
        <div className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h1 style={{ 
            margin: 0,
            background: 'linear-gradient(135deg, var(--danger) 0%, #dc2626 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent'
          }}>
            ⚙️ Admin Panel
          </h1>
          <p style={{ margin: '8px 0 0', color: 'var(--gray)', fontSize: '14px' }}>
            Manage products, inventory, and system settings
          </p>
        </div>

        {/* Create Product Section */}
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px' }}>
            ➕ Add New Product
          </h2>
          <form onSubmit={handleCreateProduct} style={{ display: 'flex', flexDirection: 'column', gap: '15px', maxWidth: '500px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                SKU (Unique Product Code) *
              </label>
              <input
                type="text"
                value={newProduct.sku}
                onChange={(e) => setNewProduct({ ...newProduct, sku: e.target.value })}
                required
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
                placeholder="e.g., PROD-001"
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Product Name *
              </label>
              <input
                type="text"
                value={newProduct.name}
                onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })}
                required
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
                placeholder="e.g., Laptop Computer"
              />
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Price ($) *
              </label>
              <input
                type="number"
                step="0.01"
                min="0"
                value={newProduct.price}
                onChange={(e) => setNewProduct({ ...newProduct, price: parseFloat(e.target.value) || 0 })}
                required
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
                placeholder="0.00"
              />
            </div>
            <button
              type="submit"
              className="btn btn-primary"
              style={{ fontSize: '16px', padding: '12px 24px' }}
            >
              ✨ Create Product
            </button>
          </form>
        </section>

        {/* Update Inventory Section */}
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px' }}>
            📦 Update Inventory
          </h2>
          <form onSubmit={handleUpdateInventory} style={{ display: 'flex', flexDirection: 'column', gap: '15px', maxWidth: '500px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                Select Product *
              </label>
              <select
                value={inventoryUpdate.product_id}
                onChange={(e) => setInventoryUpdate({ ...inventoryUpdate, product_id: e.target.value })}
                required
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
              >
                <option value="">-- Select a product --</option>
                {products.map((product) => (
                  <option key={product.id} value={product.id}>
                    {product.name} ({product.sku}) - Current: {product.inventory ?? 0}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label style={{ display: 'block', marginBottom: '5px', fontWeight: '500' }}>
                New Quantity *
              </label>
              <input
                type="number"
                min="0"
                value={inventoryUpdate.quantity}
                onChange={(e) => setInventoryUpdate({ ...inventoryUpdate, quantity: parseInt(e.target.value) || 0 })}
                required
                style={{
                  width: '100%',
                  padding: '8px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  boxSizing: 'border-box'
                }}
                placeholder="0"
              />
            </div>
            <button
              type="submit"
              className="btn btn-success"
              style={{ fontSize: '16px', padding: '12px 24px' }}
            >
              ✅ Update Inventory
            </button>
          </form>
        </section>

        {/* Products List */}
        <section className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <div className="flex-between" style={{ marginBottom: '20px' }}>
            <h2 style={{ margin: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)' }}>
              📋 All Products
            </h2>
            <button
              onClick={loadProducts}
              className="btn btn-primary"
              style={{ fontSize: '14px', padding: '10px 20px' }}
            >
              🔄 Refresh
            </button>
          </div>
          
          {loading ? (
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <div className="loading" style={{ margin: '0 auto 16px' }}></div>
              <p style={{ color: 'var(--gray)' }}>Loading products...</p>
            </div>
          ) : products.length === 0 ? (
            <div className="card" style={{ textAlign: 'center', padding: '60px 20px', background: 'rgba(99, 102, 241, 0.05)' }}>
              <div style={{ fontSize: '48px', marginBottom: '16px' }}>📦</div>
              <h3 style={{ color: 'var(--gray)', marginBottom: '8px' }}>No products found</h3>
              <p style={{ color: 'var(--gray-light)', fontSize: '14px' }}>Create your first product above!</p>
            </div>
          ) : (
            <div className="grid" style={{ gridTemplateColumns: '1fr', gap: '16px' }}>
              {products.map((product) => (
                <div
                  key={product.id}
                  className="card"
                  style={{
                    background: 'rgba(255, 255, 255, 0.8)',
                    position: 'relative',
                    overflow: 'hidden'
                  }}
                >
                  {/* Decorative gradient bar */}
                  <div style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    height: '4px',
                    background: product.inventory && product.inventory > 0 
                      ? 'linear-gradient(90deg, var(--success) 0%, var(--secondary) 100%)'
                      : 'linear-gradient(90deg, var(--danger) 0%, #dc2626 100%)'
                  }}></div>

                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginTop: '8px' }}>
                    <div style={{ flex: 1 }}>
                      <h3 style={{ marginTop: 0, marginBottom: '12px', fontSize: '1.25rem', fontWeight: '700', color: 'var(--dark)' }}>
                        {product.name}
                      </h3>
                      <div style={{ marginBottom: '12px' }}>
                        <span className="badge" style={{ 
                          background: 'rgba(99, 102, 241, 0.1)',
                          color: 'var(--primary)',
                          marginRight: '8px'
                        }}>
                          {product.sku}
                        </span>
                        <code style={{ 
                          fontSize: '0.85em', 
                          background: 'var(--gray-lighter)', 
                          padding: '4px 8px', 
                          borderRadius: '4px',
                          color: 'var(--gray)'
                        }}>
                          {product.id.slice(0, 8)}...
                        </code>
                      </div>
                      <div style={{ 
                        fontSize: '1.75rem', 
                        fontWeight: '700',
                        background: 'linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%)',
                        WebkitBackgroundClip: 'text',
                        WebkitTextFillColor: 'transparent',
                        marginBottom: '12px'
                      }}>
                        ${product.price.toFixed(2)}
                      </div>
                      {product.inventory !== undefined && (
                        <div style={{ 
                          padding: '10px',
                          borderRadius: 'var(--radius)',
                          background: product.inventory > 0 
                            ? 'rgba(16, 185, 129, 0.1)' 
                            : 'rgba(239, 68, 68, 0.1)',
                          border: `1px solid ${product.inventory > 0 ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)'}`,
                          display: 'inline-block'
                        }}>
                          <span style={{ 
                            color: product.inventory > 0 ? 'var(--success)' : 'var(--danger)',
                            fontWeight: '600',
                            fontSize: '14px'
                          }}>
                            {product.inventory > 0 ? '✅' : '❌'} {product.inventory} {product.inventory === 1 ? 'item' : 'items'} in stock
                          </span>
                        </div>
                      )}
                    </div>
                    <div style={{ display: 'flex', gap: '8px', flexDirection: 'column', marginLeft: '16px' }}>
                      <button
                        onClick={() => setInventoryUpdate({ product_id: product.id, quantity: product.inventory ?? 0 })}
                        className="btn btn-warning"
                        style={{ fontSize: '14px', padding: '8px 16px', whiteSpace: 'nowrap' }}
                      >
                        📝 Update Stock
                      </button>
                      <button
                        onClick={() => handleDeleteProduct(product.id, product.name)}
                        className="btn btn-danger"
                        style={{ fontSize: '14px', padding: '8px 16px', whiteSpace: 'nowrap' }}
                      >
                        🗑️ Delete
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
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

export default AdminPanel


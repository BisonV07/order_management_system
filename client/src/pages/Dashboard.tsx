import { useEffect, useState } from 'react'
import { productService } from '../services/api'
import type { Product } from '../types'
import '../App.css'

const Dashboard = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchProducts = async () => {
    try {
      setLoading(true)
      const data = await productService.getProducts()
      setProducts(data)
      setError(null)
    } catch (err: any) {
      const errorMsg = err?.response?.data?.message || err?.message || 'Failed to load products'
      setError(errorMsg)
      console.error('Product fetch error:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
    
    // Listen for product refresh events (triggered after order creation)
    const handleRefresh = () => {
      fetchProducts()
    }
    
    window.addEventListener('refresh-products', handleRefresh)
    
    return () => {
      window.removeEventListener('refresh-products', handleRefresh)
    }
  }, [])

  if (loading) return <div className="container">Loading...</div>
  if (error) return <div className="container">Error: {error}</div>

  return (
    <div className="app">
      <div className="container">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h1 style={{ margin: 0 }}>Product Dashboard</h1>
          <button
            onClick={fetchProducts}
            style={{
              padding: '8px 16px',
              backgroundColor: '#3498db',
              color: '#fff',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '14px',
              fontWeight: '500'
            }}
          >
            Refresh Inventory
          </button>
        </div>
        {products.length === 0 ? (
          <div style={{ padding: '20px', textAlign: 'center' }}>
            <p>No products available</p>
          </div>
        ) : (
          <div style={{ 
            display: 'grid', 
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', 
            gap: '20px',
            marginTop: '20px'
          }}>
            {products.map((product) => (
              <div 
                key={product.id}
                style={{
                  border: '1px solid #ddd',
                  borderRadius: '8px',
                  padding: '20px',
                  backgroundColor: '#fff',
                  boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
                }}
              >
                <h3 style={{ marginTop: 0, marginBottom: '10px' }}>{product.name}</h3>
                <p style={{ margin: '5px 0', color: '#666', fontSize: '0.9em' }}>
                  <strong>ID:</strong> <code style={{ backgroundColor: '#f5f5f5', padding: '2px 6px', borderRadius: '3px', fontSize: '0.85em' }}>{product.id}</code>
                </p>
                <p style={{ margin: '5px 0', color: '#666' }}>SKU: {product.sku}</p>
                <p style={{ margin: '5px 0', fontSize: '1.2em', fontWeight: 'bold', color: '#2c3e50' }}>
                  ${product.price.toFixed(2)}
                </p>
                {product.inventory !== undefined && (
                  <p style={{ 
                    margin: '5px 0', 
                    color: product.inventory > 0 ? '#27ae60' : '#e74c3c',
                    fontWeight: '500'
                  }}>
                    Stock: {product.inventory} {product.inventory === 1 ? 'item' : 'items'}
                  </p>
                )}
                {product.metadata && Object.keys(product.metadata).length > 0 && (
                  <div style={{ marginTop: '10px', paddingTop: '10px', borderTop: '1px solid #eee' }}>
                    <small style={{ color: '#999' }}>
                      {Object.entries(product.metadata).map(([key, value]) => (
                        <span key={key} style={{ marginRight: '10px' }}>
                          {key}: {String(value)}
                        </span>
                      ))}
                    </small>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Dashboard


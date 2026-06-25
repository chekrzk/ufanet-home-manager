import React, { useState } from 'react'
import { createRoot } from 'react-dom/client'
import './styles.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8086'

function App() {
  const [mode, setMode] = useState('login')
  const [tokens, setTokens] = useState(() => {
    const accessToken = localStorage.getItem('access_token')
    const refreshToken = localStorage.getItem('refresh_token')
    return accessToken ? { accessToken, refreshToken } : null
  })
  const [form, setForm] = useState({
    phone: '',
    password: '',
    full_name: '',
    house_id: '',
    apartment: '',
  })
  const [message, setMessage] = useState(null)
  const [loading, setLoading] = useState(false)

  const signedIn = Boolean(tokens?.accessToken)

  function updateField(event) {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  async function submit(event) {
    event.preventDefault()
    setMessage(null)
    setLoading(true)

    const endpoint = mode === 'login' ? '/auth/login' : '/auth/register'
    const payload =
      mode === 'login'
        ? { phone: form.phone, password: form.password }
        : {
            phone: form.phone,
            password: form.password,
            full_name: form.full_name,
            house_id: form.house_id,
            apartment: form.apartment,
          }

    try {
      const response = await fetch(`${API_URL}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      const body = await response.json().catch(() => ({}))
      if (!response.ok || body.success === false) {
        throw new Error(body.error?.message || 'Запрос не выполнен')
      }

      if (mode === 'login') {
        const data = body.data
        localStorage.setItem('access_token', data.access_token)
        localStorage.setItem('refresh_token', data.refresh_token)
        setTokens({ accessToken: data.access_token, refreshToken: data.refresh_token })
        return
      }

      setMode('login')
      setMessage({ type: 'success', text: 'Аккаунт создан. Теперь войдите.' })
    } catch (err) {
      setMessage({ type: 'error', text: err.message })
    } finally {
      setLoading(false)
    }
  }

  function changeMode(nextMode) {
    setMode(nextMode)
    setMessage(null)
  }

  function logout() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setTokens(null)
  }

  if (signedIn) {
    return (
      <main className="page">
        <section className="panel">
          <div className="badge">Личный кабинет</div>
          <h1>Вы авторизованы</h1>
          <p>Здесь будет основной интерфейс жильца: новости, заявки и уведомления.</p>
          <button className="primary" onClick={logout}>Выйти</button>
        </section>
      </main>
    )
  }

  return (
    <main className="page">
      <section className="panel">
        <div className="brand">Ufanet Home Manager</div>
        <div className="tabs">
          <button className={mode === 'login' ? 'active' : ''} onClick={() => changeMode('login')}>Вход</button>
          <button className={mode === 'register' ? 'active' : ''} onClick={() => changeMode('register')}>Регистрация</button>
        </div>

        <form onSubmit={submit}>
          <label>
            Телефон
            <input name="phone" value={form.phone} onChange={updateField} inputMode="tel" autoComplete="tel" required />
          </label>
          <label>
            Пароль
            <input name="password" type="password" value={form.password} onChange={updateField} autoComplete="current-password" required />
          </label>

          {mode === 'register' && (
            <>
              <label>
                ФИО
                <input name="full_name" value={form.full_name} onChange={updateField} autoComplete="name" required />
              </label>
              <div className="row">
                <label>
                  Дом
                  <input name="house_id" value={form.house_id} onChange={updateField} />
                </label>
                <label>
                  Квартира
                  <input name="apartment" value={form.apartment} onChange={updateField} />
                </label>
              </div>
            </>
          )}

          {message && <div className={`message ${message.type}`}>{message.text}</div>}

          <button className="primary" disabled={loading}>
            {loading ? 'Подождите...' : mode === 'login' ? 'Войти' : 'Создать аккаунт'}
          </button>
        </form>
      </section>
    </main>
  )
}

createRoot(document.getElementById('root')).render(<App />)

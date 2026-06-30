import React, { useEffect, useMemo, useState } from 'react'
import { createRoot } from 'react-dom/client'
import './styles.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8086'

const houses = [
  { id: 'house-1', name: 'Проспект Октября, 107' },
  { id: 'house-2', name: 'Улица Российская, 43' },
  { id: 'house-3', name: 'Улица Менделеева, 171' },
]

const serviceTypes = [
  { id: 'electrician', name: 'Электрик' },
  { id: 'installer', name: 'Монтажер' },
  { id: 'plumber', name: 'Сантехник' },
  { id: 'repair', name: 'Ремонтник' },
]

function decodeRole(accessToken) {
  try {
    const payload = JSON.parse(atob(accessToken.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    return payload.role || 'resident'
  } catch {
    return 'resident'
  }
}

function App() {
  const [tokens, setTokens] = useState(() => {
    const accessToken = localStorage.getItem('access_token')
    const refreshToken = localStorage.getItem('refresh_token')
    return accessToken ? { accessToken, refreshToken } : null
  })
  const [authMode, setAuthMode] = useState('login')
  const [authForm, setAuthForm] = useState({
    phone: '',
    password: '',
    full_name: '',
    role: 'resident',
  })
  const [tab, setTab] = useState('home')
  const [message, setMessage] = useState(null)
  const [profile, setProfile] = useState(null)
  const [news, setNews] = useState([])
  const [requests, setRequests] = useState([])
  const [workers, setWorkers] = useState([])
  const [availability, setAvailability] = useState([])
  const [notifications, setNotifications] = useState([])
  const [selectedHouse, setSelectedHouse] = useState('house-1')
  const [showBell, setShowBell] = useState(false)
  const [loading, setLoading] = useState(false)

  const role = useMemo(() => (tokens?.accessToken ? decodeRole(tokens.accessToken) : 'guest'), [tokens])
  const isAdmin = role === 'admin' || role === 'manager'
  const isWorker = role === 'employee'
  const houseID = isAdmin ? selectedHouse : profile?.house_id || selectedHouse

  async function api(path, options = {}) {
    const response = await fetch(`${API_URL}${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(tokens?.accessToken ? { Authorization: `Bearer ${tokens.accessToken}` } : {}),
        ...(options.headers || {}),
      },
    })
    if (response.status === 204) return null
    const body = await response.json().catch(() => ({}))
    if (!response.ok || body.success === false) {
      throw new Error(body.error?.message || 'Запрос не выполнен')
    }
    return body.data
  }

  async function loadData() {
    if (!tokens?.accessToken) return
    try {
      const [profileData, newsPage, requestsPage, notificationsPage] = await Promise.all([
        api('/profile'),
        api('/news?page=1&limit=30'),
        api('/requests?page=1&limit=30'),
        api('/notifications?page=1&limit=30'),
      ])
      setProfile(profileData)
      setNews(newsPage?.items || [])
      setRequests(requestsPage?.items || [])
      setNotifications(notificationsPage?.items || [])
    } catch (err) {
      setMessage({ type: 'error', text: err.message })
    }
  }

  async function loadHouseData(targetHouse = houseID) {
    if (!tokens?.accessToken || !targetHouse) return
    try {
      const [workerItems, availabilityItems] = await Promise.all([
        isAdmin ? api(`/profile/workers?house_id=${targetHouse}`) : Promise.resolve([]),
        api(`/profile/workers/availability?house_id=${targetHouse}`),
      ])
      setWorkers(workerItems || [])
      setAvailability(availabilityItems || [])
    } catch (err) {
      setMessage({ type: 'error', text: err.message })
    }
  }

  useEffect(() => {
    loadData()
  }, [tokens?.accessToken])

  useEffect(() => {
    loadHouseData(houseID)
  }, [tokens?.accessToken, houseID, role])

  async function submitAuth(event) {
    event.preventDefault()
    setLoading(true)
    setMessage(null)
    try {
      const endpoint = authMode === 'login' ? '/auth/login' : '/auth/register'
      const payload = authMode === 'login'
        ? { phone: authForm.phone, password: authForm.password }
        : {
            phone: authForm.phone,
            password: authForm.password,
            full_name: authForm.full_name,
            role: authForm.role,
          }
      const data = await api(endpoint, { method: 'POST', body: JSON.stringify(payload), headers: {} })
      if (authMode === 'register') {
        setAuthMode('login')
        setMessage({ type: 'success', text: 'Аккаунт создан. Теперь войдите.' })
        return
      }
      localStorage.setItem('access_token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)
      setTokens({ accessToken: data.access_token, refreshToken: data.refresh_token })
      setTab('home')
    } catch (err) {
      setMessage({ type: 'error', text: err.message })
    } finally {
      setLoading(false)
    }
  }

  async function refreshNotifications() {
    const page = await api('/notifications?page=1&limit=30')
    setNotifications(page?.items || [])
  }

  async function markRead(id) {
    await api(`/notifications/${id}/read`, { method: 'PATCH' })
    setNotifications((items) => items.map((item) => item.id === id ? { ...item, read: true } : item))
  }

  function logout() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setTokens(null)
    setProfile(null)
    setTab('home')
  }

  if (!tokens?.accessToken) {
    return <AuthScreen mode={authMode} setMode={setAuthMode} form={authForm} setForm={setAuthForm} onSubmit={submitAuth} message={message} loading={loading} />
  }

  return (
    <main className="app">
      <header className="topbar">
        <div>
          <span className="eyebrow">{roleLabel(role)}</span>
          <h1>Мой дом</h1>
        </div>
        <div className="top-actions">
          <button className="icon-button" onClick={() => { setShowBell((v) => !v); refreshNotifications() }}>🔔</button>
          <button className="icon-button" onClick={logout}>⎋</button>
        </div>
        {showBell && <NotificationsPanel items={notifications} onClose={() => setShowBell(false)} onRead={markRead} />}
      </header>

      {message && <div className={`message ${message.type}`}>{message.text}</div>}

      {tab === 'home' && <HomeScreen role={role} profile={profile} houseID={houseID} selectedHouse={selectedHouse} setSelectedHouse={setSelectedHouse} news={news} setTab={setTab} />}
      {tab === 'profile' && <ProfileScreen api={api} profile={profile} setProfile={setProfile} role={role} setMessage={setMessage} />}
      {tab === 'requests' && <RequestsScreen api={api} role={role} profile={profile} houseID={houseID} requests={requests} setRequests={setRequests} availability={availability} setAvailability={setAvailability} refreshNotifications={refreshNotifications} />}
      {tab === 'admin' && <AdminScreen api={api} selectedHouse={selectedHouse} setSelectedHouse={setSelectedHouse} workers={workers} setWorkers={setWorkers} news={news} setNews={setNews} refreshNotifications={refreshNotifications} />}

      <nav className="bottom-nav">
        <button className={tab === 'home' ? 'active' : ''} onClick={() => setTab('home')}>⌂<span>Главная</span></button>
        <button className={tab === 'requests' ? 'active' : ''} onClick={() => setTab('requests')}>☑<span>Заявки</span></button>
        <button className={tab === 'profile' ? 'active' : ''} onClick={() => setTab('profile')}>◉<span>Профиль</span></button>
        {isAdmin && <button className={tab === 'admin' ? 'active' : ''} onClick={() => setTab('admin')}>⚙<span>Админ</span></button>}
      </nav>
    </main>
  )
}

function AuthScreen({ mode, setMode, form, setForm, onSubmit, message, loading }) {
  return (
    <main className="auth-page">
      <section className="phone-frame auth-card">
        <div className="brand-mark">U</div>
        <h1>Ufanet Дом</h1>
        <p>Новости, заявки и уведомления вашего дома.</p>
        <div className="segmented">
          <button type="button" className={mode === 'login' ? 'active' : ''} onClick={() => setMode('login')}>Вход</button>
          <button type="button" className={mode === 'register' ? 'active' : ''} onClick={() => setMode('register')}>Регистрация</button>
        </div>
        <form onSubmit={onSubmit} className="form">
          <input placeholder="Телефон" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} required />
          <input placeholder="Пароль" type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} required />
          {mode === 'register' && (
            <>
              <select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
                <option value="resident">Я житель</option>
                <option value="employee">Я работник</option>
              </select>
              <input placeholder="ФИО" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} required />
            </>
          )}
          {message && <div className={`message ${message.type}`}>{message.text}</div>}
          <button className="primary" disabled={loading}>{loading ? 'Подождите...' : mode === 'login' ? 'Войти' : 'Создать аккаунт'}</button>
        </form>
      </section>
    </main>
  )
}

function HomeScreen({ role, profile, houseID, selectedHouse, setSelectedHouse, news, setTab }) {
  const visibleNews = news.filter((item) => !item.house_id || item.house_id === houseID)
  return (
    <section className="screen">
      {(role === 'admin' || role === 'manager') && (
        <label className="field">
          Дом
          <select value={selectedHouse} onChange={(e) => setSelectedHouse(e.target.value)}>
            {houses.map((house) => <option key={house.id} value={house.id}>{house.name}</option>)}
          </select>
        </label>
      )}
      <div className="hero">
        <span>{houseName(houseID)}</span>
        <h2>{profile?.full_name || 'Добро пожаловать'}</h2>
        <p>{profile?.apartment ? `Квартира ${profile.apartment}` : 'Заполните профиль и привяжитесь к дому'}</p>
      </div>
      <div className="quick-grid">
        <button onClick={() => setTab('profile')}>Профиль</button>
        <button onClick={() => setTab('requests')}>Заявки</button>
      </div>
      <h3>Лента новостей</h3>
      <div className="feed">
        {visibleNews.length === 0 && <Empty text="Пока нет новостей для выбранного дома" />}
        {visibleNews.map((item) => (
          <article className="news-card" key={item.id}>
            <span>{houseName(item.house_id)}</span>
            <h4>{item.title}</h4>
            <p>{item.body}</p>
          </article>
        ))}
      </div>
    </section>
  )
}

function ProfileScreen({ api, profile, setProfile, role, setMessage }) {
  const [form, setForm] = useState({ full_name: profile?.full_name || '', house_id: profile?.house_id || 'house-1', apartment: profile?.apartment || '' })

  async function save(event) {
    event.preventDefault()
    const data = await api('/profile', { method: 'PATCH', body: JSON.stringify(form) })
    setProfile(data)
    setMessage({ type: 'success', text: 'Профиль обновлен' })
  }

  return (
    <section className="screen">
      <h2>Профиль</h2>
      <p className="muted">Роль: {roleLabel(role)}</p>
      <form className="form" onSubmit={save}>
        <input placeholder="ФИО" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} required />
        <select value={form.house_id} onChange={(e) => setForm({ ...form, house_id: e.target.value })}>
          {houses.map((house) => <option key={house.id} value={house.id}>{house.name}</option>)}
        </select>
        {role !== 'employee' && <input placeholder="Квартира" value={form.apartment} onChange={(e) => setForm({ ...form, apartment: e.target.value })} />}
        <button className="primary">Сохранить профиль</button>
      </form>
    </section>
  )
}

function RequestsScreen({ api, role, profile, houseID, requests, setRequests, availability, setAvailability, refreshNotifications }) {
  const [form, setForm] = useState({ category: 'electrician', preferred_date: '', assigned_worker_id: '', description: '', address: houseName(houseID), apartment: profile?.apartment || '', phone: '' })
  const [workTime, setWorkTime] = useState({ specialization: 'electrician', available_date: '', available_time: '' })
  const filteredAvailability = availability.filter((item) => item.specialization === form.category && (!form.preferred_date || item.available_date === form.preferred_date))

  async function createRequest(event) {
    event.preventDefault()
    const item = await api('/requests', { method: 'POST', body: JSON.stringify(form) })
    setRequests([item, ...requests])
    await refreshNotifications()
  }

  async function saveAvailability(event) {
    event.preventDefault()
    const item = await api('/profile/workers/availability', {
      method: 'POST',
      body: JSON.stringify({ ...workTime, house_id: houseID }),
    })
    setAvailability([item, ...availability])
  }

  async function setStatus(request, status) {
    const updated = await api(`/requests/${request.id}/status`, { method: 'PATCH', body: JSON.stringify({ status, assigned_to: request.assigned_to || profile?.id || '' }) })
    setRequests(requests.map((item) => item.id === updated.id ? updated : item))
    await refreshNotifications()
  }

  return (
    <section className="screen">
      <h2>Заявки</h2>
      {role === 'employee' ? (
        <form className="form compact" onSubmit={saveAvailability}>
          <select value={workTime.specialization} onChange={(e) => setWorkTime({ ...workTime, specialization: e.target.value })}>
            {serviceTypes.map((type) => <option key={type.id} value={type.id}>{type.name}</option>)}
          </select>
          <input type="date" value={workTime.available_date} onChange={(e) => setWorkTime({ ...workTime, available_date: e.target.value })} required />
          <input type="time" value={workTime.available_time} onChange={(e) => setWorkTime({ ...workTime, available_time: e.target.value })} required />
          <button className="secondary">Указать доступность</button>
        </form>
      ) : (
        <form className="form compact" onSubmit={createRequest}>
          <select value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value, assigned_worker_id: '' })}>
            {serviceTypes.map((type) => <option key={type.id} value={type.id}>{type.name}</option>)}
          </select>
          <input type="date" value={form.preferred_date} onChange={(e) => setForm({ ...form, preferred_date: e.target.value, assigned_worker_id: '' })} />
          <select value={form.assigned_worker_id} onChange={(e) => setForm({ ...form, assigned_worker_id: e.target.value })}>
            <option value="">Работник будет назначен</option>
            {filteredAvailability.map((item) => <option key={item.id} value={item.user_id}>{serviceName(item.specialization)} · {item.available_date} {item.available_time}</option>)}
          </select>
          <input placeholder="Адрес" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
          <input placeholder="Квартира" value={form.apartment} onChange={(e) => setForm({ ...form, apartment: e.target.value })} />
          <input placeholder="Телефон для связи" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
          <textarea placeholder="Опишите проблему" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} required />
          <button className="primary">Отправить заявку</button>
        </form>
      )}
      <div className="cards">
        {requests.length === 0 && <Empty text="Заявок пока нет" />}
        {requests.map((request) => (
          <article className="request-card" key={request.id}>
            <div><strong>{serviceName(request.category)}</strong><span>{request.status}</span></div>
            <p>{request.description}</p>
            <small>{request.address || houseName(houseID)} · кв. {request.apartment || '-'} · {request.phone || 'без телефона'}</small>
            {role === 'employee' && request.status !== 'done' && (
              <div className="button-row">
                <button className="secondary" onClick={() => setStatus(request, 'in_progress')}>Да</button>
                <button className="secondary" onClick={() => setStatus(request, 'canceled')}>Нет</button>
                <button className="primary" onClick={() => setStatus(request, 'done')}>Завершено</button>
              </div>
            )}
          </article>
        ))}
      </div>
    </section>
  )
}

function AdminScreen({ api, selectedHouse, setSelectedHouse, workers, setWorkers, news, setNews, refreshNotifications }) {
  const [newsForm, setNewsForm] = useState({ title: '', body: '' })
  const [workerForm, setWorkerForm] = useState({ full_name: '', specialization: 'electrician', phone: '', user_id: '' })

  async function createNews(event) {
    event.preventDefault()
    const item = await api('/news', { method: 'POST', body: JSON.stringify({ ...newsForm, house_id: selectedHouse }) })
    setNews([item, ...news])
    setNewsForm({ title: '', body: '' })
    await refreshNotifications()
  }

  async function addWorker(event) {
    event.preventDefault()
    const worker = await api('/profile/workers', { method: 'POST', body: JSON.stringify({ ...workerForm, house_id: selectedHouse }) })
    setWorkers([worker, ...workers])
    setWorkerForm({ full_name: '', specialization: 'electrician', phone: '', user_id: '' })
  }

  return (
    <section className="screen">
      <h2>Панель управления</h2>
      <label className="field">
        Выбор дома
        <select value={selectedHouse} onChange={(e) => setSelectedHouse(e.target.value)}>
          {houses.map((house) => <option key={house.id} value={house.id}>{house.name}</option>)}
        </select>
      </label>
      <form className="form compact" onSubmit={createNews}>
        <h3>Создать новость</h3>
        <input placeholder="Заголовок" value={newsForm.title} onChange={(e) => setNewsForm({ ...newsForm, title: e.target.value })} required />
        <textarea placeholder="Текст новости" value={newsForm.body} onChange={(e) => setNewsForm({ ...newsForm, body: e.target.value })} required />
        <button className="primary">Опубликовать</button>
      </form>
      <form className="form compact" onSubmit={addWorker}>
        <h3>Работники</h3>
        <select value={workerForm.specialization} onChange={(e) => setWorkerForm({ ...workerForm, specialization: e.target.value })}>
          {serviceTypes.map((type) => <option key={type.id} value={type.id}>{type.name}</option>)}
        </select>
        <input placeholder="ФИО работника" value={workerForm.full_name} onChange={(e) => setWorkerForm({ ...workerForm, full_name: e.target.value })} required />
        <input placeholder="Телефон" value={workerForm.phone} onChange={(e) => setWorkerForm({ ...workerForm, phone: e.target.value })} />
        <input placeholder="User ID аккаунта работника" value={workerForm.user_id} onChange={(e) => setWorkerForm({ ...workerForm, user_id: e.target.value })} />
        <button className="secondary">Добавить работника</button>
      </form>
      <div className="service-groups">
        {serviceTypes.map((type) => (
          <div className="group" key={type.id}>
            <h4>{type.name}</h4>
            {workers.filter((worker) => worker.specialization === type.id).map((worker) => <p key={worker.id}>{worker.full_name} · {worker.phone || 'без телефона'}</p>)}
          </div>
        ))}
      </div>
    </section>
  )
}

function NotificationsPanel({ items, onClose, onRead }) {
  return (
    <div className="notifications">
      <div><strong>Уведомления</strong><button onClick={onClose}>×</button></div>
      {items.length === 0 && <p className="muted">Пока нет уведомлений</p>}
      {items.map((item) => (
        <article key={item.id} className={item.read ? 'read' : ''} onClick={() => !item.read && onRead(item.id)}>
          <b>{item.title}</b>
          <p>{item.body}</p>
        </article>
      ))}
    </div>
  )
}

function Empty({ text }) {
  return <div className="empty">{text}</div>
}

function houseName(id) {
  return houses.find((house) => house.id === id)?.name || id || 'Дом не выбран'
}

function serviceName(id) {
  return serviceTypes.find((type) => type.id === id)?.name || id
}

function roleLabel(role) {
  return ({ admin: 'Администратор', manager: 'Управляющий', employee: 'Работник', resident: 'Житель', guest: 'Гость' })[role] || role
}

createRoot(document.getElementById('root')).render(<App />)

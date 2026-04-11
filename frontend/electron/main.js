const { app, BrowserWindow, ipcMain, Menu, shell, Tray, nativeImage } = require('electron')
const path = require('path')
const childProcess = require('child_process')
const net = require('net')
const treeKill = require('tree-kill')

let mainWindow = null
let tray = null
let backendChild = null
app.isQuitting = false

// ─── Backend Process ───────────────────────────────────────────
function isPortInUse(port) {
  return new Promise((resolve) => {
    const client = new net.Socket()
    client.once('error', () => {
      client.destroy()
      resolve(false)
    })
    client.once('connect', () => {
      client.destroy()
      resolve(true)
    })
    client.connect(port, '127.0.0.1')
  })
}

function killBackend() {
  if (backendChild && backendChild.pid) {
    treeKill(backendChild.pid, 'SIGTERM', (err) => {
      if (err) {
        console.error('[CATendar] Failed to kill backend:', err.message)
      } else {
        console.log('[CATendar] Backend killed')
      }
    })
  }
}

async function startBackend() {
  const portInUse = await isPortInUse(18082)
  if (portInUse) {
    console.log('[CATendar] Backend already running on port 18082')
    return
  }

  const exeName = 'calendar-backend.exe'
  const exePath = app.isPackaged
    ? path.join(process.resourcesPath, exeName)
    : path.join(__dirname, '../build', exeName)

  const workDir = app.isPackaged
    ? path.dirname(path.join(process.resourcesPath, exeName))
    : path.join(__dirname, '../build')

  console.log('[CATendar] Starting backend from:', exePath)
  backendChild = childProcess.spawn(exePath, [], {
    stdio: 'ignore',
    cwd: workDir,
    env: { ...process.env, DB_PATH: path.join(__dirname, '../../backend/calendar.db') },
    windowsHide: true,
  })

  // Wait for backend to actually be ready (max 10s)
  const maxWait = 10000
  const interval = 300
  let waited = 0
  while (!(await isPortInUse(18082))) {
    if (waited >= maxWait) {
      console.error('[CATendar] Backend failed to start within 10s')
      return
    }
    await new Promise((r) => setTimeout(r, interval))
    waited += interval
  }
  console.log('[CATendar] Backend ready on port 18082')
}

// ─── Tray Icon ─────────────────────────────────────────────────
function createTray() {
  const iconPath = path.join(__dirname, '../dist/assets/CATendar.ico')
  let icon = nativeImage.createFromPath(iconPath)
  icon = icon.resize({ width: 16, height: 16 })
  tray = new Tray(icon)

  tray.setToolTip('Calendar')

  tray.on('click', () => {
    if (mainWindow) {
      if (mainWindow.isVisible()) {
        mainWindow.hide()
      } else {
        mainWindow.show()
        mainWindow.focus()
      }
    }
  })

  tray.setContextMenu(createTrayMenu())
}

function createTrayMenu() {
  return Menu.buildFromTemplate([
    {
      label: '显示窗口',
      click: () => {
        mainWindow.show()
        mainWindow.focus()
      }
    },
    { type: 'separator' },
    {
      label: '退出',
      click: () => {
        app.isQuitting = true
        killBackend()
        app.quit()
      }
    }
  ])
}

// ─── Main Window ──────────────────────────────────────────────
function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1000,
    height: 700,
    minWidth: 360,
    minHeight: 280,
    icon: path.join(__dirname, '../dist/assets/CATendar.ico'),
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false
    },
    frame: true,
    titleBarStyle: 'default'
  })

  const isDev = !app.isPackaged

  if (isDev) {
    mainWindow.loadURL('http://localhost:5179')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
  }

  mainWindow.on('close', (event) => {
    if (!app.isQuitting) {
      event.preventDefault()
      mainWindow.hide()
    }
  })

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

// ─── App Lifecycle ─────────────────────────────────────────────
app.whenReady().then(async () => {
  const isDev = !app.isPackaged

  if (!isDev) {
    app.setLoginItemSettings({ openAtLogin: true })
    await startBackend()
  } else {
    console.log('[CATendar] Dev mode — backend managed by concurrently')
  }

  Menu.setApplicationMenu(null)
  createWindow()
  createTray()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  // Don't quit — stay in tray
})

app.on('before-quit', () => {
  killBackend()
})

// ─── IPC Handlers ──────────────────────────────────────────────
ipcMain.handle('set-always-on-top', (_, flag, level) => {
  if (mainWindow) {
    mainWindow.setAlwaysOnTop(flag, level || 'normal')
  }
})

ipcMain.handle('open-external', (_, url) => {
  if (url) {
    shell.openExternal(url)
  }
})

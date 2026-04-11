// Starts the backend binary and waits for it to be ready
const { spawn } = require('child_process')
const http = require('http')
const path = require('path')

const exePath = path.join(__dirname, '../build/calendar-backend.exe')
const dbPath = path.join(__dirname, '../../backend/calendar.db')

console.log('[start-backend] Launching:', exePath)
const child = spawn(exePath, [], {
  stdio: 'ignore',
  cwd: path.join(__dirname, '../build'),
  env: { ...process.env, DB_PATH: dbPath },
  windowsHide: true,
  detached: true,
})
child.unref()

function waitForBackend() {
  return new Promise((resolve, reject) => {
    const maxWait = 30000
    const interval = 500
    let waited = 0

    const check = () => {
      const req = http.get('http://127.0.0.1:18082/api/health', (res) => {
        if (res.statusCode === 200) {
          console.log('[start-backend] Backend ready!')
          resolve()
        } else {
          retry()
        }
      })
      req.on('error', retry)
      req.setTimeout(2000, () => {
        req.destroy()
        retry()
      })
    }

    const retry = () => {
      if (waited >= maxWait) {
        reject(new Error('[start-backend] Backend failed to start within 30s'))
      } else {
        waited += interval
        setTimeout(check, interval)
      }
    }

    console.log('[start-backend] Waiting for backend on port 18082...')
    check()
  })
}

waitForBackend()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err.message)
    process.exit(1)
  })

const { contextBridge, ipcRenderer } = require('electron')

contextBridge.exposeInMainWorld('electronAPI', {
  setAlwaysOnTop: (flag, level = 'normal') => ipcRenderer.invoke('set-always-on-top', flag, level),
  openExternal: (url) => ipcRenderer.invoke('open-external', url)
})

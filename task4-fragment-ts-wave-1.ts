// ВНИМАНИЕ: в этом фрагменте есть несколько ошибок и плохих практик.
// Кандидату нужно:
// 1) Найти и описать проблемы.
// 2) Предложить, как переписать код лучше.

type Device = {
  id: number
  hostname: string
  ip: string
}

// Имитация запроса к API
async function fetchDevices(): Promise<Device[]> {
  // Потенциальная проблема: игнорируются ошибки сети/HTTP-код
  const res = await fetch('/api/devices')

  if(!res.ok) {
    throw new Error(res.statusText)
  }

  return (await res.json()) as Device[]
}


function normalizeString(str: string): string {
  return str.trim().toLocaleLowerCase()
}

class DeviceService {
  private devices: Device[]
  private isLoading: boolean
  private errors: Error[]

  constructor() {
    this.devices = []
    this.isLoading = false
    this.errors = []
  }

  async loadAndFilterDevices(search: string): Promise<Device[]> {
    try {
      this.isLoading = true

  const data = await fetchDevices()

  // Потенциальная проблема: мутируем общий массив из разных мест
      this.devices = data

  // Потенциальная проблема: сравнение без нормализации регистра и trim
      const filtered = this.devices.filter((d) => normalizeString(d.hostname).indexOf(normalizeString(search)) >= 0)

  // Потенциальная проблема: функция имеет побочные эффекты и возвращает разные типы
  // (в реальном коде сюда часто добавляют ещё логику, что делает её трудной для тестирования)
  return filtered
    } catch(e) {
      console.error(e)
      this.errors.push(e as Error)
    } finally {
      this.isLoading = false
    }
    
    return []
  }
}
}

// Пример использования (упрощённо)
async function example() {
  const searchInput: HTMLInputElement | null = document.querySelector('#search')
  if (searchInput) {
    // Потенциальная проблема: нет debounce, каждый ввод символа может бить по API
    searchInput.oninput = async () => {
      const list = await loadAndFilterDevices(searchInput.value)
      console.log('Devices:', list)
    }
  }
}

example()


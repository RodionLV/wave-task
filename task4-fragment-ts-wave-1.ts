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

// Глобальное состояние (антипаттерн для большинства приложений)
let devices: Device[] = []
let isLoading = false

export async function loadAndFilterDevices(search: string) {
  isLoading = true

  // Потенциальная проблема: нет try/catch, при ошибке состояние "подвиснет"
  const data = await fetchDevices()

  // Потенциальная проблема: мутируем общий массив из разных мест
  devices = data

  // Потенциальная проблема: сравнение без нормализации регистра и trim
  const filtered = devices.filter((d) => d.hostname.indexOf(search) >= 0)

  // "забываем" сбросить isLoading в случае ошибок выше
  isLoading = false

  // Потенциальная проблема: функция имеет побочные эффекты и возвращает разные типы
  // (в реальном коде сюда часто добавляют ещё логику, что делает её трудной для тестирования)
  return filtered
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


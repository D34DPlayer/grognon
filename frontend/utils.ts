export function displayTime(time: Date | string | number): string {
  if (!time) {
    return ''
  }
  if (!(time instanceof Date)) {
    switch (typeof time) {
      case 'string':
        time = new Date(time)
        break
      case 'number':
        time = new Date(time * 1000) // Assuming the number is in seconds
        break
      default:
        return ''
    }
  }

  return time.toLocaleString()
}

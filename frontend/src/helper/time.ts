export type Timeline = { start: number; steps: number; range: number; stepSize: number }

export const calculateStart = (range: number, stepSize: number): number => {
  const start = (new Date()).getTime() - range

  return start - (start % stepSize)
}

export const timeline = (time: number|null, steps = 12): Timeline => {
  const range = (time || 5 * 60000)
  const stepSize = range / steps

  return {
    range,
    steps,
    start: calculateStart(range, stepSize),
    stepSize
  }
}

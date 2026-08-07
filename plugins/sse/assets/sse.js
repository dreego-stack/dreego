export function connectSSE(url, onEvent) {
  const es = new EventSource(url)
  es.onmessage = (e) => onEvent(e.data)
  return es
}

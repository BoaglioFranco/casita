export function formatAddedTime(createdAt: string | undefined, now = Date.now()): string {
  if (!createdAt) return 'Added time unavailable';

  const timestamp = Date.parse(createdAt);
  if (!Number.isFinite(timestamp)) return 'Added time unavailable';

  const age = Math.max(0, now - timestamp);
  if (age < 60_000) return 'just now';
  if (age < 86_400_000) {
    const unit = age < 3_600_000 ? 'minute' : 'hour';
    const amount = Math.floor(age / (unit === 'minute' ? 60_000 : 3_600_000));
    return new Intl.RelativeTimeFormat(undefined, { numeric: 'always' }).format(-amount, unit);
  }

  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  }).format(timestamp);
}

import {
  normalizeDate,
  toValue,
  useNow,
  type DateLike,
  type MaybeRefOrGetter,
} from '@vueuse/core';
import { computed } from 'vue';

function formatHours(hours: number) {
  if (hours === 1) {
    return '1 hour ago';
  }
  return `${hours} hours ago`;
}

function formatMinutes(minutes: number) {
  if (minutes === 1) {
    return '1 minute ago';
  }
  return `${minutes} minutes ago`;
}

export const formatRelativeDate = function (
  date: Date,
  now: Date = new Date(),
) {
  const diff = now.getTime() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days === 0) {
    if (hours === 0) {
      if (minutes === 0) {
        return 'Just now';
      }
      return formatMinutes(minutes);
    }
    return formatHours(hours);
  }
  if (days === 1) {
    return `Yesterday at ${date.toLocaleTimeString()}`;
  }
  return date.toLocaleString();
};

export const useRelativeFormattedDate = function (
  date: MaybeRefOrGetter<DateLike>,
) {
  const now = useNow({ interval: 1000 });
  return computed(() =>
    formatRelativeDate(normalizeDate(toValue(date)), toValue(now)),
  );
};

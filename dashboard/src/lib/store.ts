import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const localStore = (key: string, initial: string) => {
  if (browser && localStorage.getItem(key) === null) {
    localStorage.setItem(key, initial);
  }


  const saved = browser ? localStorage.getItem(key) : "";

  const { subscribe, set, update } = writable(saved);

  return {
    subscribe,
    set: (value: string) => {
      localStorage.setItem(key, value);
      return set(value);
    },
    update,
  };
};

export const store = localStore('api-key', '');
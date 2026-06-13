export const colors = {
  background: '#FAFAF7',
  primary: '#111111',
  secondary: '#666666',
  border: '#ECECEC',
  white: '#FFFFFF',
} as const;

export type ColorName = keyof typeof colors;

export const moodColors: Record<string, string> = {
  still: '#A8B5C8',
  waiting: '#9AA3B2',
  drift: '#B9C7C0',
  warm: '#D8B48C',
  distant: '#A8A8B8',
  quiet: '#B5B5A8',
  hollow: '#C2C2C2',
  soft: '#D4C4B0',
  fading: '#B8B8B8',
  wondering: '#A9B0C0',
  returning: '#B8C4A8',
};

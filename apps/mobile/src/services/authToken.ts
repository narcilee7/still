export type GetToken = () => Promise<string | null>;

let getTokenImpl: GetToken = async () => null;

export function setGetToken(fn: GetToken) {
  getTokenImpl = fn;
}

export function getToken(): Promise<string | null> {
  return getTokenImpl();
}

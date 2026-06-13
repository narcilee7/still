import React, { useEffect } from 'react';
import { useAuth } from '@clerk/expo';
import { setGetToken } from '../services/authToken';

/**
 * TokenBridge wires Clerk's session token into the API transport.
 * Render it once inside <ClerkProvider>.
 */
export function TokenBridge() {
  const { getToken } = useAuth();

  useEffect(() => {
    setGetToken(() => getToken());
  }, [getToken]);

  return null;
}

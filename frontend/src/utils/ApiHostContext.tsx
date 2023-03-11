import { createContext } from 'react';

type ApiHost = {
  host: string;
  port: number;
};

export function DefaultApiHost(): ApiHost {
  return {
    host: process.env.REACT_APP_APIHOST_NAME!,
    port: Number(process.env.REACT_APP_APIHOST_PORT!),
  };
}

export const ApiHostContext = createContext<ApiHost>(DefaultApiHost());

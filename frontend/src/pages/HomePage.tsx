import { Button, Card, Typography } from '@mui/material';
import React, { useContext } from 'react';
import { ApiHostContext, DefaultApiHost } from '../utils/ApiHostContext';

const HomePage: React.FC = () => {
  return (
    <ApiHostContext.Provider value={DefaultApiHost()}>
      <div>
        <h1>トップページ</h1>
        <img
          src={`${process.env.PUBLIC_URL}/logo512.png`}
          alt='ra-noyokushinryu'
        />
      </div>
      <ApiCallTestComponent />
    </ApiHostContext.Provider>
  );
};

const ApiCallTestComponent: React.FC = () => {
  const [numUsers, setNumUsers] = React.useState<number | undefined>(undefined);
  const apiHostCtx = useContext(ApiHostContext);

  return (
    <Card variant='outlined'>
      <Typography sx={{ fontSize: 14 }}>
        numUsers: {numUsers === undefined ? 'undefined' : numUsers}
      </Typography>
      <Button
        onClick={() => {
          fetch(`http://${apiHostCtx.host}:${apiHostCtx.port}/api/hello`, {
            method: 'POST',
            body: JSON.stringify({
              param1: 'Ra',
              param2: 100,
            }),
            headers: {
              'Content-type': 'application/json; charset=UTF-8',
            },
          })
            .then((response) => response.json())
            .then((data) => {
              console.log(data.message);
              setNumUsers(data.numUsers as number);
            })
            .catch((err) => {
              console.log(err.message);
            });
        }}
      >
        Hello
      </Button>
    </Card>
  );
};

export default HomePage;

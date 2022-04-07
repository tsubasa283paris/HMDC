import React from 'react';

import Router from './Router';
import NavBar from './components/NavBar';

const App: React.FC = () => {
  return (
    <React.Fragment>
      <NavBar />
      <Router />
    </React.Fragment>
  );
};

export default App;
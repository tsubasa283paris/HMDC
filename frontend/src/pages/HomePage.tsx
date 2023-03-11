import React from "react";

const HomePage: React.FC = () => {
  return (
    <div>
      <h1>トップページ</h1>
      <img src={`${process.env.PUBLIC_URL}/logo512.png`} alt="ra-noyokushinryu"/>
    </div>
  );
};

export default HomePage;
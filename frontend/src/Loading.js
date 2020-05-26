import React from "react";
import Lottie from "lottie-react-web";

import connecting from "./assets/loading.json";

export default function Loading() {
  return (
    <div className="loading">
      <Lottie
        width="200px"
        height="200px"
        options={{
          animationData: connecting
        }}
      />
      <div className="loading_text">Cooking up ideas...</div>
    </div>
  );
}

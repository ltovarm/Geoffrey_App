import React, { useState, useEffect } from "react";

function App() {
  const [temperatures, setTemperatures] = useState([]);

  useEffect(() => {
    // Connect to WebSocket
    const socket = new WebSocket("ws://localhost:8080/ws");

    // Listen for messages from the server
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setTemperatures(data.temperatures);
    };

    // Close the WebSocket connection when unmounting
    return () => {
      socket.close();
    };
  }, []);

  return (
    <div>
      <h1>Temperaturas</h1>
      {temperatures.map((temp, index) => (
        <p key={index}>Temperatura {index + 1}: {temp.toFixed(1)} °C</p>
      ))}
    </div>
  );
}

export default App;

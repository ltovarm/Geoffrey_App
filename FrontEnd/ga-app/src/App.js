import React, { useState, useEffect } from "react";

function App() {
  const [temperatures, setTemperatures] = useState([]);
  const [buttonColors, setButtonColors] = useState(Array(5).fill("red"));

  useEffect(() => {
    // Obtener las temperaturas desde el backend en Go
    fetchTemperatures()
      .then((data) => setTemperatures(data.temperatures))
      .catch((error) => console.error("Error al obtener las temperaturas:", error));
  }, []);

  const fetchTemperatures = async () => {
    const response = await fetch("http://localhost:8080/temperatures");
    const data = await response.json();
    return data;
  };

  useEffect(() => {
    // Guardar las temperaturas en el almacenamiento local del navegador
    localStorage.setItem("temperatures", JSON.stringify(temperatures));
  }, [temperatures]);

  const handleButtonClick = (index) => {
    // Cambiar el color del botón al ser pulsado
    const newButtonColors = [...buttonColors];
    newButtonColors[index] = newButtonColors[index] === "red" ? "green" : "red";
    setButtonColors(newButtonColors);
  };

  return (
    <div>
      <h1>Temperaturas</h1>
      {temperatures.map((temp, index) => (
        <p key={index}>Temperatura {index + 1}: {temp.toFixed(2)} °C</p>
      ))}
      <div>
        {buttonColors.map((color, index) => (
          <button
            key={index}
            style={{ backgroundColor: color }}
            onClick={() => handleButtonClick(index)}
          >
            Botón {index + 1}
          </button>
        ))}
      </div>
    </div>
  );
}

export default App;

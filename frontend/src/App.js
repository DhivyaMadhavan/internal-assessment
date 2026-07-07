import "./App.css";
import { useEffect, useState } from "react";

import Header from "./components/Header";
import CircuitPanel from "./components/CircuitPanel";
import MetricsPanel from "./components/MetricsPanel";
import StatsCards from "./components/StatsCards";
import TransitionLog from "./components/TransitionLog";

function App() {

    const [data, setData] = useState({
        requestCount: 0,
        rps: 0,
        route: "Primary",
        circuitState: "CLOSED"
    });

    const [transitionLog, setTransitionLog] = useState([]);
    const [graphData, setGraphData] = useState([]);

    useEffect(() => {

        const ws = new WebSocket("ws://localhost:8080/ws");

        ws.onmessage = (event) => {
            console.log("Received:", event.data);

            const msg = JSON.parse(event.data);
            const safeRps = Number(msg.rps ?? 0);

            const timestamp = Date.now();

            // MAIN STATE UPDATE
            setData(prev => {

                if (prev.circuitState !== msg.circuitState) {

                    const now = new Date().toLocaleTimeString();

                    setTransitionLog(old => [
                        `${now} : ${prev.circuitState} → ${msg.circuitState}`,
                        ...old.slice(0, 9)
                    ]);
                }

                return {
                    ...msg,
                    rps: safeRps
                };
            });

           

            // GRAPH UPDATE (more stable for Recharts tooltip)
            setGraphData(old => {

                const updated = [
                    ...old,
                    {
                        time: timestamp,   // 🔥 important fix
                        rps: safeRps
                    }
                ];

                return updated.slice(-20);
            });

        };

        ws.onclose = () => {
            console.log("WebSocket disconnected");
        };

        return () => ws.close();

    }, []);

    useEffect(() => {
                    console.log("graphData:", graphData);
                }, [graphData]);

    return (
        <div className="container">

            <Header />

            <div className="topRow">

                <CircuitPanel data={data} />

                <MetricsPanel
                    graphData={graphData}
                    currentRps={data.rps}
                />

            </div>

            <StatsCards data={data} />

            <TransitionLog logs={transitionLog} />

        </div>
    );
}

export default App;
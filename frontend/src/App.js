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
        route: "Primary",
        circuitState: "CLOSED"
    });

    const [transitionLog, setTransitionLog] = useState([]);
    const [graphData, setGraphData] = useState([]);

    useEffect(() => {

        const ws = new WebSocket("ws://localhost:8080/ws");

        ws.onmessage = (event) => {

            const msg = JSON.parse(event.data);

            // Update transition log and latest metrics
            setData(prev => {

                if (prev.circuitState !== msg.circuitState) {

                    const now = new Date().toLocaleTimeString();

                    setTransitionLog(old => [

                        `${now} : ${prev.circuitState} → ${msg.circuitState}`,

                        ...old.slice(0, 9)

                    ]);

                }

                return msg;

            });

            // Update graph
            const now = new Date().toLocaleTimeString();

            setGraphData(old => {

                const updated = [
                    ...old,
                    {
                        time: now,
                        requests: msg.requestCount
                    }
                ];

                // Keep only last 20 data points
                if (updated.length > 20) {
                    updated.shift();
                }

                return updated;

            });

        };

        ws.onclose = () => {
            console.log("WebSocket disconnected");
        };

        return () => ws.close();

    }, []);

    return (

        <div className="container">

            <Header />

            <div className="topRow">

                <CircuitPanel data={data} />

                <MetricsPanel graphData={graphData} />

            </div>

            <StatsCards data={data} />

            <TransitionLog logs={transitionLog} />

        </div>

    );

}

export default App;
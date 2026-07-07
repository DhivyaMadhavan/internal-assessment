import RequestsChart from "./RequestsChart";

function MetricsPanel({ graphData, currentRps }) {

    return (

        <div className="panel">

            <h3>REQUESTS / SEC</h3>

            {/* 🔥 FIX 4: show numeric value clearly */}
            <div style={{ fontSize: "20px", marginBottom: "10px" }}>
                Current RPS: <b>{currentRps}</b>
            </div>

            <RequestsChart data={graphData} />

        </div>

    );

}

export default MetricsPanel;
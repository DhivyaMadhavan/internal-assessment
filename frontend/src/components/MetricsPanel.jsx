import RequestsChart from "./RequestsChart";

function MetricsPanel({ graphData }) {

    return (

        <div className="panel">

            <h3>REQUESTS / SEC</h3>

            <RequestsChart data={graphData} />

        </div>

    );

}

export default MetricsPanel;
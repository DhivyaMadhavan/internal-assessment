import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    CartesianGrid
} from "recharts";

function CustomTooltip({ active, payload }) {
    if (!active || !payload?.length) return null;

    const point = payload[0].payload;

    return (
        <div style={{
            background: "#111",
            color: "#fff",
            padding: "8px 10px",
            borderRadius: "6px",
            fontSize: "12px"
        }}>
            <div><b>Time:</b> {new Date(point.time).toLocaleTimeString()}</div>
            <div><b>RPS:</b> {point.rps}</div>
        </div>
    );
}

function RequestsChart({ data }) {

    return (
        <ResponsiveContainer width="100%" height={250}>

            <LineChart data={data}>

                <CartesianGrid strokeDasharray="3 3" />

                <XAxis
                    dataKey="time"
                    tickFormatter={(t) =>
                        new Date(t).toLocaleTimeString()
                    }
                />

                <YAxis />

                <Tooltip content={<CustomTooltip />} />

                {/* 🔥 THIS MUST MATCH BACKEND + APP.JS */}
                <Line
                    type="monotone"
                    dataKey="rps"
                    stroke="#8884d8"
                    dot={false}
                    isAnimationActive={false}
                />

            </LineChart>

        </ResponsiveContainer>
    );
}

export default RequestsChart;
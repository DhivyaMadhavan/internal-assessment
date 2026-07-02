import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer
} from "recharts";

function RequestsChart({ data }) {

    return (

        <ResponsiveContainer width="100%" height={250}>

            <LineChart data={data}>

                <XAxis dataKey="time" />

                <YAxis />

                <Tooltip />

                <Line
                    type="monotone"
                    dataKey="requests"
                    stroke="#00ff88"
                    strokeWidth={2}
                    dot={false}
                />

            </LineChart>

        </ResponsiveContainer>

    );

}

export default RequestsChart;
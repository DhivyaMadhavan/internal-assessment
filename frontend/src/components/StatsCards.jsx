function StatsCards({data}){

    return(

        <div className="statsGrid">

            <div className="statCard">
                <small>Total Requests</small>
                <h2>{data.requestCount}</h2>
            </div>

            <div className="statCard">
                <small>Route</small>
                <h2>{data.route}</h2>
            </div>

            <div className="statCard">
                <small>Circuit</small>
                <h2>{data.circuitState}</h2>
            </div>

        </div>

    )

}

export default StatsCards;
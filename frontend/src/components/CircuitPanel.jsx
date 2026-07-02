function CircuitPanel({data}){

    return(

        <div className="panel">

            <h3>CIRCUIT BREAKER</h3>

            <div className="diagram">

                <div className="box">
                    Client
                </div>

                <div className="arrow">→</div>

                <div className="box">
                    Router
                </div>

                <div className="arrow">→</div>

                <div className={
                    data.route==="Primary"
                    ?"box active"
                    :"box"
                }>
                    Primary
                </div>

            </div>

            <div className="diagram">

                <div style={{width:"170px"}}></div>

                <div className="arrow">↘</div>

                <div className={
                    data.route==="Secondary"
                    ?"box activeFallback"
                    :"box"
                }>
                    Secondary
                </div>

            </div>

            <h2 className={data.circuitState.toLowerCase()}>
                {data.circuitState}
            </h2>

        </div>

    )

}

export default CircuitPanel;
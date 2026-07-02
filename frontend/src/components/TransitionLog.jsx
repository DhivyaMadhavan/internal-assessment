function TransitionLog({ logs }) {

    return (

        <div className="panel">

            <h3>TRANSITION LOG</h3>

            {
                logs.length === 0 ?

                    <p>No state transitions yet.</p>

                    :

                    logs.map((log, index) => (

                        <p key={index}>{log}</p>

                    ))
            }

        </div>

    );

}

export default TransitionLog;
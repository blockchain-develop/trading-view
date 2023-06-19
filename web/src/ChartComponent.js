import { useEffect, useRef } from "react"
import { createChart } from 'lightweight-charts';
import { priceData } from "./priceData";

function ChartComponent() {

    const chartContainerRef = useRef();

    useEffect(
        () => {
            const handleResize = () => {
                chart.applyOptions({width: chartContainerRef.current.clientWithd});
            };

            //
            const chart = createChart(chartContainerRef.current, {
                layout: {
                    background: {type:'solid'.Solid, color:'white'},
                    textColor:'black',
                },
                width: chartContainerRef.current.clientWithd,
                height: 300,
            });
            chart.timeScale().fitContent();

            const priceSeries = chart.addCandlestickSeries({
                upColor: "#4bffb5",
                downColor: "#ff4976",
                borderDownColor: "#ff4976",
                borderUpColor: "#4bffb5",
                wickDownColor: "#838ca1",
                wickUpColor: "#838ca1",
                priceLineVisible: false,
            });

            //
            //
            var priceData1;
            fetch('http://127.0.0.1:8000/api/price', {
                method: 'POST',
                mode: 'cors',
                headers: {
                    'Accept':'application/json',
                    'Content-Type':'application/json'
                },
                body:JSON.stringify({})
            })
            .then(response => response.json())
            .then(json => {
                console.log(json)
                priceData1 = json;
                priceSeries.setData(priceData1);
            })

            //
            //priceSeries.setData(priceData);
            
            window.addEventListener('resize', handleResize);

            return () => {
                window.removeEventListener('resize', handleResize);
                chart.remove();
            };
        },
    );

    return (
        <div
            ref={chartContainerRef}
            className="chart-container"
        />
    );
};

export default ChartComponent
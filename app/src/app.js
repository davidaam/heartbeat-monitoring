import moment from "moment"
import * as d3 from "d3";
import TimelineChart from "d3-timeline-chart"
import request from "request"

request('http://127.0.0.1:8080/heartbeats', (error, response, body) => {
  !error || console.log('Error:', error);
  let json = JSON.parse(body).reduce((json, datapoint) => {
    json[datapoint.clientID] = json[datapoint.clientID] || []
    json[datapoint.clientID].push(datapoint.timestamp)
    return json
  }, {})
  let data = Object.keys(json).map((clientID) => {
    json[clientID].sort((x, y) => x - y)
    let intervals =
      json[clientID].slice(1, json[clientID].length -1).reduce((intervals, timestamp) => {
        let currInterval = intervals[intervals.length-1]
        if (currInterval.to && timestamp - currInterval.to > 1) {
          intervals.push({ type: TimelineChart.TYPE.INTERVAL, from: timestamp, to: timestamp, })
        } else {
          currInterval.to = timestamp
        }
        return intervals
      }, [{
        from: json[clientID][0],
        to: json[clientID][0]
      }])

    return {
      label: clientID,
      data: intervals.map((interval) => ({
        type: TimelineChart.TYPE.INTERVAL,
        from: moment(interval.from*1000).toDate(),
        to: moment(interval.to*1000).toDate()
      }))
    }
  })

  const element = document.getElementById('heartbeat-chart');
  console.log(data)
  const timeline = new TimelineChart(element, data, {
      tip: function(d) {
          return `${d.from}<br>${d.to}`;
      }
  });
})

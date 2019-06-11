package main

const (
	tmplDebugReport = `<!DOCTYPE html>
<html>
  <head>
    <title>Embedding Vega-Lite</title>
    <script src="https://cdn.jsdelivr.net/npm/vega@5.4.0"></script>
    <script src="https://cdn.jsdelivr.net/npm/vega-lite@3.3.0"></script>
    <script src="https://cdn.jsdelivr.net/npm/vega-embed@4.2.0"></script>
  </head>
  <body>
    <div id="cumulative"></div>
    <div id="timeline"></div>

    <script type="text/javascript">
      var cumulativeSpec = {
        $schema: 'https://vega.github.io/schema/vega-lite/v2.0.json',
        width: 1200,
        height: 400,
        description: 'cumulative distribute',
        data: {
          values: %s
        },
        mark: 'bar',
        encoding: {
          x: {field: 'v', type: 'ordinal'},
          y: {field: 'c', type: 'quantitative'}
        }
      };
	  vegaEmbed('#cumulative', cumulativeSpec);

      var timelineSpec = {
        $schema: "https://vega.github.io/schema/vega-lite/v2.json",
        width: 1200,
        height: 200,
        data: {values: %s},
        mark: {
          type: 'circle',
          opacity: 0.8,
          stroke: 'black',
          strokeWidth: 1,
        },
        encoding: {
          x: {field: 'date', type: 'temporal', timeUnit:"hoursminutesseconds", scale: {type: 'utc'}},
          y: {field: 'tid', type: 'nominal', axis: {title: ""}},
          size: {field: 'c', type: 'quantitative', legend: {title: 'Sample Counts', clipHeight: 30}},
          color: {field: 'tid', type: 'nominal', legend: null},
        }
      };
      vegaEmbed('#timeline', timelineSpec);
    </script>
  </body>
</html>`
)

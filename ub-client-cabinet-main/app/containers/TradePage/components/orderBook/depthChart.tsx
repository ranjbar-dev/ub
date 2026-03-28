// import React, { memo, useState, useLayoutEffect } from 'react';
// import { Sparklines, SparklinesLine } from 'react-sparklines';
// import { DepthChartSubscriber } from 'services/message_service';

// interface Props {
//   color: string;
//   side: 'depthBuy' | 'depthSell';
// }

// function DepthChart(props: Props) {
//   const { color, side } = props;
//   const [ChartData, setChartData] = useState([0]);
//   useLayoutEffect(() => {
//     const DepthChartSubscription = DepthChartSubscriber.subscribe(
//       (message: any) => {
//         setChartData(message[side]);
//       },
//     );
//     return () => {
//       DepthChartSubscription.unsubscribe();
//     };
//   }, []);
//   return (
//     <Sparklines data={ChartData}>
//       <SparklinesLine color={color} style={{ strokeWidth: 1 }} />
//     </Sparklines>
//   );
// }

// export default memo(DepthChart);

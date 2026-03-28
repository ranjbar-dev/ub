//import * as d3 from 'd3';
import {scaleLinear,extent,min,max,select,area,line} from 'd3';

export const ResizeGridHeigth=(data: {
	uniqueId: string;
	additinal?: number;
}) => {
	const layOutWrapper=document.getElementById(data.uniqueId+'Container');
	if(layOutWrapper) {
		const layoutHeight=layOutWrapper.style.height;
		const gridWrapper=document.getElementById(
			'ag-grid-wrapper-'+data.uniqueId,
		);
		if(gridWrapper) {
			gridWrapper.style.height=
				+layoutHeight.replace('px','')-
				40-
				(data.additinal? data.additinal:0)+
				'px';
		}
	}
};
export const LayoutWidth=() => {
	return window.innerWidth;
};
export const LayoutHeight=() => {
	const availableHeight=window.innerHeight-60;
	if(availableHeight>600) {
		return +availableHeight/177.5;
	}
	return 5;
};
export const layoutMargin=() => {
	const availableHeight=window.innerHeight-60;

	const verticalMargin=+(availableHeight/88.7).toFixed(2);
	if(availableHeight>600) {
		return [10,verticalMargin];
	}
	return [10,10];
};
export const chartParams: any={
	halfWidth: 0,
	chartHeight: 0,
	x: 0,
	y: 0,
	greenData: 0,
	redData: 0,
	Area: 0,
	svg: 0,
};
export const DrawDepthChart=async (data: {
	sellData: any[];
	buyData: any[];
	update?: boolean;
}) => {
	const mainContainer=document.getElementById('ORDERBOOKContainer');
	if(mainContainer) {
		chartParams.halfWidth=mainContainer.clientWidth/2-4;
		chartParams.chartHeight=
			+mainContainer.style.height.replace('px','')-55;
	}
	//////////////////
	//range
	chartParams.x=scaleLinear().range([0,chartParams.halfWidth]);
	chartParams.y=scaleLinear().range([chartParams.chartHeight,0]);

	//area
	chartParams.greenData=data.buyData;
	chartParams.redData=data.sellData;

	DepthSvgRenderer(
		chartParams.x,
		chartParams.chartHeight,
		chartParams.y,
		chartParams.greenData,
		chartParams.halfWidth+1,
		'green',
		data.update,
	);
	DepthSvgRenderer(
		chartParams.x,
		chartParams.chartHeight,
		chartParams.y,
		chartParams.redData,
		chartParams.halfWidth+1,
		'red',
		data.update,
	);
	return;
	/////////////////////////
};
export const colorAndInexes=(type) => {
	return type==='green'
		? [
			{
				backgroundColor: 'rgba(6, 186, 97, 0.18)',
				offset: 30,
				opacity: 1,
			},
			{
				backgroundColor: 'rgba(70, 214, 143, 0.12)',
				offset: 60,
				opacity: 1,
			},
		]
		:[
			{
				backgroundColor: 'rgba(244, 6, 70, 0.19)',
				offset: 30,
				opacity: 1,
			},
			{
				backgroundColor: 'rgba(244, 6, 70, 0.07)',
				offset: 60,
				opacity: 1,
			},
		];
};
function DepthSvgRenderer(
	x: any,
	chartHeight: any,
	y: any,
	chartData: any[],
	halfWidth: any,
	type: string,
	update?: boolean,
) {
	if(update!==true) {
		// define the area
		chartParams.Area=area()
			.x(function(d) {
				return x(d['price']);
			})
			.y0(chartHeight)
			.y1(function(d) {
				return y(+d['sum']);
			});
		// Scale the range of the data
		x.domain(
			extent(chartData,function(d) {
				return +d['price'];
			}) as any,
		);
		y.domain([
			min(chartData,function(d) {
				return +d['sum'];
			}) as number,
			max(chartData,function(d) {
				return +d['sum'];
			}) as number,
		]);

		// append the svg object to the body of the container
		// appends a 'group' element to 'svg'
		chartParams.svg=select('.orderBookchart'+(type==='green'? '1':'2'))
			.html('')
			.append('svg')
			.attr('class','svgorderBookchart'+(type==='green'? '1':'2'))
			.attr('width',type==='green'? halfWidth-2:halfWidth+4)
			.attr('height',chartHeight)
			.attr(
				'style',
				'margin-top:0px;margin-left:'+
				(type==='green'? -2:halfWidth-3)+
				'px;',
			);
		//  .append('g');
		//  .attr('transform', 'translate(0,0)')
		//  .append('g');
		// .call(axis);

		//chartParams.svg;
		// set the gradient

		chartParams.svg
			.append('defs')
			.append('linearGradient')
			.attr('gradientUnits','userSpaceOnUse')
			.attr(
				'id',
				'gradient-areaorderBookchart'+(type==='green'? '1':'2'),
			)
			.attr('x1',0)
			.attr('y1',0)
			.attr('x2',0)
			.attr('y2',y(0)||0)
			.selectAll('stop')
			.data(colorAndInexes(type)) // <= your new offsets array
			.enter()
			.append('stop')
			.attr('offset',function(d) {
				return d.offset+'%';
			})
			.attr('style',function(d) {
				return 'stop-color:'+d.backgroundColor+';stop-opacity:'+d.opacity;
			});
		// Add the area.
		chartParams.svg
			.append('path')
			.data([chartData])
			.attr('class','area')
			.attr(
				'style',
				'fill: url(#gradient-areaorderBookchart'+
				(type==='green'? '1':'2')+
				');',
			)
			.attr('stroke','none')
			.attr('d',chartParams.Area);
		//add line chart
		chartParams.svg
			.append('path')
			.datum(chartData)
			.attr('fill','none')
			.attr('stroke',type==='green'? '#69b3a2':'rgb(244, 6, 70)')
			.attr('stroke-width',1)
			.attr(
				'd',
				line()
					.x(function(d) {
						return x(d['price']);
					})
					.y(function(d) {
						return y(d['sum']);
					}) as any,
			);
	} else {
	}
}
// export const DepthSvgDrawer

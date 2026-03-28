import ReactGA from 'react-ga'
const GA_TRACKING_ID='UA-211160738-1'

export const initGA=() => {
	ReactGA.initialize(GA_TRACKING_ID)
	const path=window.location.pathname==='/'? 'app-home':window.location.pathname
	ReactGA.set({page: path})
	ReactGA.pageview(path)
}

export const logPageView=(url) => {
	ReactGA.set({page: url})
	ReactGA.pageview(url)
}

export const logEvent=(category='',action='') => {
	if(category&&action) {
		ReactGA.event({category,action})
	}
}

export const logException=(description='',fatal=false) => {
	if(description) {
		ReactGA.exception({description,fatal})
	}
}
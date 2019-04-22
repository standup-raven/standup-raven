function getStyle() {
    return {
        header: {
            marginTop: '0',
            marginBottom: '30px',
        },
        controlBtns: {
            marginRight: '10px',
            width: '65px',
            height: '30px',
        },
        form: {
            height: '250px',
            overflowY: 'auto',
            marginBottom: '16.5px',
        },
        alert: {
            width: '90%',
            marginLeft: 'auto',
            marginRight: 'auto',
            textAlign: 'center',
            borderRadius: '5px',
            whiteSpace: 'pre-line',
            animation: 'pop 0.2s ease-in',
        },
        spinner: {
            width: '80px',
            display: 'block',
            margin: '50px auto',
        },
        standupErrorMessage: {
            fontWeight: 'bold',
        },
    };
}

export default {
    getStyle,
};

const computeMonthlyOnTheWhich = (data, rruleObj) => {
    if (rruleObj.freq !== 1 || !rruleObj.bysetpos) {
        return data.repeat.monthly.onThe.which;
    }

    const bysetpos = (typeof rruleObj.bysetpos === 'number') ? rruleObj.bysetpos : rruleObj.bysetpos[0];

    switch (bysetpos) {
        case 1: {
            return 'First';
        }
        case 2: {
            return 'Second';
        }
        case 3: {
            return 'Third';
        }
        case 4: {
            return 'Fourth';
        }
        case -1: {
            return 'Last';
        }
        default: {
            return data.repeat.monthly.onThe.which;
        }
    }
};

export default computeMonthlyOnTheWhich;

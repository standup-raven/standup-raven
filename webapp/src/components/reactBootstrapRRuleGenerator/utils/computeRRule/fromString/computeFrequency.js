const computeFrequency = (data, rruleObj) => {
    switch (rruleObj.freq) {
        case 0: {
            return 'Yearly';
        }
        case 1: {
            return 'Monthly';
        }
        case 2: {
            return 'Weekly';
        }
        case 3: {
            return 'Daily';
        }
        case 4: {
            return 'Hourly';
        }
        default: {
            return data.repeat.frequency;
        }
    }
};

export default computeFrequency;

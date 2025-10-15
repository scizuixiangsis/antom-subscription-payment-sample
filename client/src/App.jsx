import './App.css';
import { useEffect, useState } from 'react';

const paymentOptions = [
    {
        label: 'Standard',
        value: 'Standard',
        icon: 'https://cdn.marmot-cloud.com/storage/2022/12/20/4e873d78-0f2d-4553-895e-e5392eddda39.png',
        amountValue: '688.80',
        currency: 'HKD',
        onboardingTime: '3-5 working days',
        languageSupported: 'Chinese/English',
        inquiryResponseTime: '< 48 hours',
    },
    {
        label: 'Premium',
        value: 'Premium',
        icon: '	https://cdn.marmot-cloud.com/storage/2022/12/19/d45c0b0f-a1f3-4ca9-99c7-f0ea74c5c6ac.png',
        amountValue: '1088.80',
        currency: 'HKD',
        onboardingTime: '1-2 working days',
        languageSupported: 'Chinese/English/Japanese/Korean',
        inquiryResponseTime: '< 24 hours',
    },
];

function App() {
    const initialPaymentOption =
        paymentOptions.find(
            (option) =>
                option.value ===
                (localStorage.getItem('supportPlan') || 'Standard')
        ) || paymentOptions[0];

    const [selectedOption, setSelectedOption] = useState(initialPaymentOption);
    const [isBtnLoading, setBtnLoading] = useState(false);
    const [paymentResult, setPaymentResult] = useState({});
    const [alert, setAlert] = useState({
        title: '',
        subTitle: '',
        iconType: '',
    });

    const handlePaymentOptionChange = (value) => {
        setPaymentResult({});
        const option = paymentOptions.find((option) => option.value === value);
        if (option) {
            setSelectedOption(option);
            localStorage.setItem('supportPlan', option.value);
        }
    };

    // Page error handling
    const renderAlert = () => {
        if (!alert.title) return;

        const icons = {
            warning:
                'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*7ShHRKjYcjoAAAAAAAAAAAAAAQAAAQ',
            error: 'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*QzBHQ40jKZAAAAAAAAAAAAAAAQAAAQ',
            success:
                'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*ZLggSbkLKoMAAAAAAAAAAAAAAQAAAQ',
        };

        return (
            <>
                <div className="container-desktop desktop-animation">
                    <div className="container-content">
                        <img src={icons[alert.iconType]} alt="" />
                        <div className="container-title">{alert.title}</div>
                        <div className="container-subTitle">
                            {alert.subTitle}
                        </div>
                        <button
                            type="button"
                            className="container-btn"
                            onClick={() =>
                                setAlert({
                                    title: '',
                                    subTitle: '',
                                    iconType: '',
                                })
                            }
                        >
                            OK
                        </button>
                    </div>
                </div>
                <div className="mockup"></div>
            </>
        );
    };

    // Step 4:Display the payment results at the bottom of the page.
    const renderPaymentResult = () => {
        if (!Object.keys(paymentResult).length) return null;

        const { title, message, icon, bgColor, borderColor } = paymentResult;

        return (
            <div
                className="result-container"
                style={{ backgroundColor: bgColor, borderColor: borderColor }}
            >
                <div className="result-title">
                    <img src={icon} alt="" id="resultIcon" />
                    <span id="message">{title}</span>
                </div>
                <div id="content">{message}</div>
            </div>
        );
    };

    // Step 3: Select the rendered data according to the result of the payment.
    const handlePaymentStatus = (status) => {
        setTimeout(() => {
            window.location.href = './index.html';
            localStorage.removeItem('supportPlan');
        }, 2000);

        switch (status) {
            case 'SUCCESS':
                setPaymentResult({
                    title: 'Payment Successful',
                    message:
                        'Thank you for your payment! We will ship out your order as soon as possible.',
                    icon: 'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*ZLggSbkLKoMAAAAAAAAAAAAAAQAAAQ',
                    bgColor: '#e5f7f1',
                    borderColor: '#b7e9d9',
                });
                break;
            case 'FAIL':
                setPaymentResult({
                    title: 'Payment Failed',
                    message:
                        'Please return to the merchant order page and resubmit your payment.',
                    icon: 'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*QzBHQ40jKZAAAAAAAAAAAAAAAQAAAQ',
                    bgColor: 'rgba(255, 91, 77, 0.10)',
                    borderColor: 'rgba(255, 91, 77, 0.20)',
                });
                break;
            case 'ERROR':
                setPaymentResult({
                    title: 'Error',
                    message:
                        'An error occurred while checking the payment status. Please try again or contact support.',
                    icon: 'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*QzBHQ40jKZAAAAAAAAAAAAAAAQAAAQ',
                    bgColor: 'rgba(255, 91, 77, 0.10)',
                    borderColor: 'rgba(255, 91, 77, 0.20)',
                });
                break;
            default:
                setPaymentResult({
                    title: 'Payment Processing',
                    message:
                        'Please return to the merchant order page and resubmit your payment.',
                    icon: 'https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*7ShHRKjYcjoAAAAAAAAAAAAAAQAAAQ',
                    bgColor: 'rgba(255, 159, 26, 0.10)',
                    borderColor: 'rgba(255, 159, 26, 0.20) ',
                });
        }
    };
    // Step 2: Get Subscription Payment Result
    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search);
        const resultStatus = searchParams.get('resultStatus');
        if (resultStatus) {
            handlePaymentStatus(resultStatus);
        }
    }, []);
    // Step 1: Get consult data from the endpoint.
    const handleSubmit = async () => {
        setPaymentResult({});
        setBtnLoading(true);

        if (!selectedOption.value) {
            setAlert({
                title: 'Warning',
                subTitle: 'Please select a payment method!',
                iconType: 'warning',
            });
            setBtnLoading(false);
            return null;
        }

        try {
            const url = 'http://localhost:8080/subscriptions/create';

            const body = JSON.stringify({
                periodCount: '4',
                periodType: 'YEAR',
                currency: selectedOption.currency,
                amountValue: selectedOption.amountValue,
                paymentMethodType: 'ALIPAY_HK',
            });

            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body,
            });

            if (!response.ok) {
                throw new Error('HTTP error! Status: ' + response.status);
            }
            const res = await response.json();
            if (
                res.status === 'success' &&
                res.data.result.resultStatus === 'S' &&
                res.data.normalUrl
            ) {
                location.href = res.data.normalUrl;
            }        } catch (error) {
            setAlert({
                title: 'Error',
                subTitle:
                    error.message || 'Failed to fetch payment session data',
                iconType: 'error',
            });
            return null;
        } finally {
            setBtnLoading(false);
        }
    };
    return (
        <>
            <div className="antom-content">
                <header>
                    <div className="header-content">
                        <img
                            src="https://mdn.alipayobjects.com/portal_pdqp4x/afts/file/A*xWIWSb_-O6MAAAAAAAAAAAAAAQAAAQ"
                            alt=""
                            className="header-content-img"
                        />
                        <div className="subtitle">
                            Subscription Payment Demo
                        </div>
                    </div>
                </header>
                <div className="content">
                    <div className="checkout-container">
                        <div className="container-left">
                            <div className="leftContent">
                                <div className="detail-row">
                                    <div className="detail-header">
                                        Support plan
                                    </div>
                                    <div className="detail-content">
                                        {selectedOption.label}
                                    </div>
                                </div>
                                <div className="detail-row">
                                    <div className="detail-header">
                                        Onboarding time
                                    </div>
                                    <div className="detail-content">
                                        {selectedOption.onboardingTime}
                                    </div>
                                </div>
                                <div className="detail-row">
                                    <div className="detail-header">
                                        Language supported
                                    </div>
                                    <div className="detail-content">
                                        {selectedOption.languageSupported}
                                    </div>
                                </div>
                                <div className="detail-row">
                                    <div className="detail-header">
                                        Inquiry response time
                                    </div>
                                    <div className="detail-content">
                                        {selectedOption.inquiryResponseTime}
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className="container-right">
                            <div className="rightContent">
                                <div className="product-info">
                                    <p>
                                        Select the support plan you want to
                                        subscribe to
                                    </p>
                                    <div
                                        className="payment-options radio-group"
                                        id="checkout-box"
                                    >
                                        {paymentOptions.map(
                                            ({ label, value }) => {
                                                return (
                                                    <label                                                        className={`payment-option ${selectedOption.value === value ? 'selected' : ''}`}                                                        key={value}
                                                    >
                                                        <input
                                                            type="radio"
                                                            className="radio-input"
                                                            name="radio"
                                                            value={value}
                                                            onChange={(
                                                                event
                                                            ) => {
                                                                handlePaymentOptionChange(
                                                                    event.target
                                                                        .value
                                                                );
                                                            }}
                                                            checked={
                                                                selectedOption.value ===
                                                                value
                                                            }
                                                        />
                                                        <span className="radio-button"></span>
                                                        {label}
                                                    </label>
                                                );
                                            }
                                        )}
                                    </div>
                                    <button
                                        className="submit"
                                        id="submit"
                                        onClick={handleSubmit}
                                        disabled={isBtnLoading}
                                    >
                                        Pay {selectedOption.currency}{' '}
                                        {selectedOption.amountValue}
                                    </button>
                                </div>

                                {renderPaymentResult()}
                            </div>
                        </div>
                    </div>
                    {renderAlert()}
                </div>
                <div className="footer">
                    This page is used only for sandbox testing, for more testing
                    information, please refer
                    <a
                        href="https://global.alipay.com/docs/ac/cashierpay/test"
                        className="gotoTest"
                        target="_blank"
                    >
                        Testing resources
                    </a>
                </div>
            </div>
        </>
    );
}
export default App;

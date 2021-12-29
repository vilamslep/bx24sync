USE dev

DECLARE @order BINARY(16)
DECLARE @yearsOffset NUMERIC(4)
DECLARE @docType VARCHAR(20)

SELECT @yearsOffset = offset FROM dbo._YearOffset

SET @order = ${order}
SET @docType = 'Заказ клиента'

SELECT 
	_bx24User._Fld18234RRef as oneCUser,
	_bx24User._Fld18232 as bx24User
INTO #bx24Users
FROM dbo._Reference18202_VT18230 _bx24User
WHERE NOT _bx24User._Fld18234RRef = 0x00000000000000000000000000000000

SELECT
	_order._IDRRef as ref,
	LOWER( 
		CONCAT(
			SUBSTRING(CONVERT(VARCHAR(40),_order._IDRRef,2),25,8), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._IDRRef,2),21,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._IDRRef,2),17,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._IDRRef,2),1,4), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._IDRRef,2),5 , 12)
			)
	) AS originId,
	FORMAT(_order._Date_Time, 'dd.MM.yyyy', 'ru-RU') as docDate,
	CONCAT(@docType, ' ', _order._Number) as name,
	LOWER(
		CONCAT(
			SUBSTRING(CONVERT(VARCHAR(40), _order._Fld7829RRef,2),25,8), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._Fld7829RRef,2),21,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._Fld7829RRef,2),17,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._Fld7829RRef,2),1,4), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_order._Fld7829RRef,2),5 , 12)
			)
	) client,
	_order._Fld7835 as docSum,
	_internetOrderStage._EnumOrder+1 as internetOrderStage,
	_order._Fld7892 as dpp,
	_order._Fld7893 as dppOD,
	_order._Fld7894 as dppOS,
	_pickUpPoint._Description as pickupPoint,
	_order._Fld17567 as internetOrder,
	_order._Fld17569 as sentSms,
	_orderType._EnumOrder+1 as orderType,
	_order._Fld18175 as deliverySum,
	_deliveryWay._EnumOrder+1 as deliveryWay,
	FORMAT(_order._Fld7837, 'dd.MM.yyyy', 'ru-RU') as wantedDateShipment,
	REPLACE( REPLACE( CAST(_order._Fld7882 AS VARCHAR(1024)), ';',' ') , char(10), ' ') as extraInfo,
	REPLACE( REPLACE( CAST(_order._Fld7860 AS VARCHAR(1024)), ';',' ') , char(10), ' ') as comment,
	_agreement._Description as agreement,
	_stock._Description as stock,
	REPLACE(CAST(_order._Fld7856 AS VARCHAR(1024)), ';',' ') as deliveryAddress,
	ISNULL(_deliveryArea._Description, '') as deliveryArea,
	FORMAT(_order._Fld7877, 'dd.MM.yyyyThh:mm:ss', 'ru-RU') as deliveryTimeFrom,
	FORMAT(_order._Fld7878, 'dd.MM.yyyyThh:mm:ss', 'ru-RU') as deliveryTimeTo,
	ISNULL(_person._Description, '') as doctor,
	CASE WHEN _bx24User.bx24User IS NULL THEN 475 ELSE _bx24User.bx24User END as userId,
	FORMAT(_order._Fld7855, 'dd.MM.yyyy', 'ru-RU') as shipmentDate,
	_departmentParent._Description as department
INTO #orders
FROM dbo._Document367 _order WITH(NOLOCK)
LEFT OUTER JOIN dbo._Enum17526 _internetOrderStage WITH(NOLOCK)
ON _order._Fld17528RRef = _internetOrderStage._IDRRef
LEFT OUTER JOIN dbo._Reference17525 _pickUpPoint WITH(NOLOCK)
ON (_order._Fld17530RRef = _pickUpPoint._IDRRef)
LEFT OUTER JOIN dbo._Enum17897 _orderType WITH(NOLOCK)
ON _order._Fld17898RRef = _orderType._IDRRef
LEFT OUTER JOIN dbo._Enum651 _deliveryWay WITH(NOLOCK)
ON _order._Fld7874RRef = _deliveryWay._IDRRef
LEFT OUTER JOIN dbo._Reference231 _agreement WITH(NOLOCK)
ON (_order._Fld7832RRef = _agreement._IDRRef)
LEFT OUTER JOIN dbo._Reference297 _stock WITH(NOLOCK)
ON (_order._Fld7838RRef = _stock._IDRRef)
LEFT OUTER JOIN dbo._Reference102 _deliveryArea WITH(NOLOCK)
ON (_order._Fld7876RRef = _deliveryArea._IDRRef)
LEFT OUTER JOIN dbo._Reference265 _person
ON (_order._Fld7891RRef = _person._IDRRef)
LEFT OUTER JOIN #bx24Users as _bx24User WITH(NOLOCK)
ON (_order._Fld7840RRef = _bx24User.oneCUser)
LEFT OUTER JOIN dbo._Reference244 _department WITH(NOLOCK)
ON (_order._Fld7870RRef = _department._IDRRef)
LEFT OUTER JOIN dbo._Reference244 _departmentParent WITH(NOLOCK)
ON (_department._ParentIDRRef = _departmentParent._IDRRef)
WHERE _order._IDRRef = @order
SELECT
	_paymentSchedules._Document367_IDRRef as ref,
	_paymentSchedules._Fld7928RRef as paymentOption,
	_paymentSchedules._Fld7931 as paymentSum
INTO #paymentSchedules
FROM #orders _order WITH(NOLOCK)
LEFT OUTER JOIN dbo._Document367_VT7926 _paymentSchedules WITH(NOLOCK)
ON ((_order.ref = _paymentSchedules._Document367_IDRRef)) 

SELECT
	t1.ref,
	t1.paymentSum as prepaid,
	CAST(0.0 AS NUMERIC(15, 2)) as prepayment,
	CAST(0.0 AS NUMERIC(15, 2)) as credit
INTO #paymentByKind
FROM #paymentSchedules t1 WITH(NOLOCK)
WHERE (t1.paymentOption= 0x94C1BD2F840B2F7C4145EA35E822F82C)

UNION ALL 

SELECT
	t2.ref,
	CAST(0.0 AS NUMERIC(15, 2)),
	t2.paymentSum,
	CAST(0.0 AS NUMERIC(15, 2))
FROM #paymentSchedules t2 WITH(NOLOCK)
WHERE (t2.paymentOption= 0xBEF8B570E82E62D5458EA7D2753EDF13)

UNION ALL 

SELECT
	t3.ref,
	CAST(0.0 AS NUMERIC(15, 2)),
	CAST(0.0 AS NUMERIC(15, 2)),
	t3.paymentSum
FROM #paymentSchedules t3 WITH(NOLOCK)
WHERE (t3.paymentOption= 0x8B7EC90DA0FFA2DF42842BE394A0DE53)

SELECT
	t1.ref as ref,
	CAST(SUM(t1.prepaid) AS NUMERIC(15, 2)) as prepaid,
	CAST(SUM(t1.prepayment) AS NUMERIC(15, 2)) as prepayment,
	CAST(SUM(t1.credit) AS NUMERIC(15, 2)) as credit
INTO #schedules
FROM #paymentByKind  t1 WITH(NOLOCK)
GROUP BY t1.ref

SELECT 
	t1.*,
	ISNULL(t2.prepaid, 0) as prepaid,
	ISNULL(t2.prepayment, 0) as prepayment,
	ISNULL(t2.credit, 0) as credit
FROM
	#orders t1
	LEFT OUTER JOIN #schedules t2
		ON t1.ref = t2.ref

DROP TABLE #bx24Users
DROP TABLE #schedules
DROP TABLE #paymentByKind
DROP TABLE #paymentSchedules
DROP TABLE #orders
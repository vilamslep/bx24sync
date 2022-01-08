USE dev

DECLARE @shipment BINARY(16)
DECLARE @emptyRef BINARY(16)
DECLARE @clinic BINARY(16)
DECLARE @docType VARCHAR(20)
SET @shipment = ${shipment}
SET @docType = 'Реализация'
SET @emptyRef = 0x00000000000000000000000000000000
SET @clinic = 0x98AC00025553E65311E3631E1F2C2507

SELECT 
	_bx24User._Fld18234RRef as oneCUser,
	_bx24User._Fld18232 as bx24User
INTO #bx24Users
FROM dbo._Reference18202_VT18230 _bx24User
WHERE NOT _bx24User._Fld18234RRef = @emptyRef

SELECT
	CONVERT(VARCHAR(40), t1._IDRRef,2) AS ref,
	LOWER(CONCAT(
			SUBSTRING(CONVERT(VARCHAR(40), t1._IDRRef,2),25,8), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),t1._IDRRef,2),21,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),t1._IDRRef,2),17,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),t1._IDRRef,2),1,4), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),t1._IDRRef,2),5 , 12)
	)) as originId,
	CONCAT(@docType, ' ', t1._Number) as name,
	FORMAT(t1._Date_Time, 'dd.MM.yyyy', 'ru-RU') as docDate,
	CONVERT(VARCHAR(40),t1._Fld10612RRef,2) as client,
	t1._Fld10611 as docSum,
	t4._Description as department,
	t5._Description as stock,
	t6._Description as agreement,
	REPLACE(
		REPLACE(
			CAST(t1._Fld10622 AS VARCHAR(1024)), ';',' '
		)
		, char(10), ' '
	) as comment,
	ISNULL(t9._Description, '') as doctor,
	CASE WHEN t2.bx24User IS NULL THEN 475 ELSE t2.bx24User END as userId
FROM dbo._Document431 t1 WITH(NOLOCK)
LEFT OUTER JOIN #bx24Users t2 WITH(NOLOCK)
ON t1._Fld10609RRef = t2.oneCUser
LEFT OUTER JOIN dbo._Reference244 t3 WITH(NOLOCK)
ON t1._Fld10614RRef = t3._IDRRef
LEFT OUTER JOIN dbo._Reference244 t4 WITH(NOLOCK)
ON t3._ParentIDRRef = t4._IDRRef
LEFT OUTER JOIN dbo._Reference297 t5 WITH(NOLOCK)
ON t1._Fld10617RRef = t5._IDRRef
LEFT OUTER JOIN dbo._Reference231 t6 WITH(NOLOCK)
ON t1._Fld10619RRef = t6._IDRRef
LEFT OUTER JOIN dbo._Reference102 t8 WITH(NOLOCK)
ON t1._Fld10638RRef = t8._IDRRef
LEFT OUTER JOIN dbo._Reference265 t9 WITH(NOLOCK)
ON t1._Fld17149RRef = t9._IDRRef
WHERE 
	t1._Posted = 0x01 
	AND (t1._Fld10606_TYPE = 0x08 AND t1._Fld10606_RTRef = 0x0000016F AND t1._Fld10606_RRRef = @emptyRef) 
	AND T1._Fld10607RRef <> @clinic
	AND t1._IDRRef = @shipment

DROP TABLE #bx24Users
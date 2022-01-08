USE dev

DECLARE @reception BINARY(16)
DECLARE @yearsOffset NUMERIC(4)

SET @reception = ${reception}
SELECT @yearsOffset = offset FROM dbo._YearOffset

SELECT
	CONVERT(VARCHAR(40),_reception._IDRRef,2) as id, 
	LOWER(
		CONCAT(
			SUBSTRING(CONVERT(VARCHAR(40),_reception._IDRRef,2),25,8), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_reception._IDRRef,2),21,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_reception._IDRRef,2),17,4),
			'-', SUBSTRING(CONVERT(VARCHAR(40),_reception._IDRRef,2),1,4), 
			'-', SUBSTRING(CONVERT(VARCHAR(40),_reception._IDRRef,2),5 , 12)
			)
	) as originId,
	CONCAT(_reception._Fld17546, ' ',_reception._Number) name,
	CASE
			WHEN _reception._Date_Time = '2001-01-01 00:00:00.000'
			THEN NULL
			ELSE FORMAT(DATEADD(YEAR, -@yearsOffset, _reception._Date_Time), 'dd.MM.yyyy', 'ru-RU')
		END as date,
	CONVERT(VARCHAR(40),_reception._Fld12347RRef,2) client,
	_departmentParent._Description department,
	_bx24usr._Fld18232 as userId
FROM dbo._Document484 _reception WITH(NOLOCK)
	LEFT OUTER JOIN dbo._Reference244 _department WITH(NOLOCK)
	ON (_reception._Fld17250RRef = _department._IDRRef)
	LEFT OUTER JOIN dbo._Reference244 _departmentParent WITH(NOLOCK)
	ON (_department._ParentIDRRef = _departmentParent._IDRRef)
	LEFT OUTER JOIN dbo._Reference18202_VT18230 _bx24usr WITH(NOLOCK)
	LEFT OUTER JOIN dbo._Reference291 _users WITH(NOLOCK)
	ON (_bx24usr._Fld18234RRef = _users._IDRRef)
	ON ((_reception._Fld12348RRef = _users._Fld6504RRef))
WHERE _reception._IDRRef = @reception


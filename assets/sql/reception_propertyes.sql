USE dev

DECLARE @reception BINARY(16)
DECLARE @yearsOffset NUMERIC(4)

SET @reception = ${reception}
SELECT @yearsOffset = offset FROM dbo._YearOffset

SELECT DISTINCT
	_propety._Fld18365 as id,
	CONVERT(int, _params._Fld15106_TYPE, 2) as valType,
	_params._Fld15106_L as valBool,
	_params._Fld15106_N as valDigit,
	CASE
			WHEN _params._Fld15106_T = '2001-01-01 00:00:00.000'
			THEN NULL
			ELSE FORMAT(DATEADD(YEAR, -@yearsOffset, _params._Fld15106_T), 'dd.MM.yyyy', 'ru-RU')
		END AS valTime,
	_params._Fld15106_S as valStr,
	CASE
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 107 THEN t1._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 123 THEN t2._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 265 THEN t3._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 286 THEN t4._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 295 THEN t5._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 17405 THEN t6._Description
		WHEN CONVERT(int, _params._Fld15106_RTRef) = 17871 THEN t7._Description
	END as valRef
FROM dbo._Chrc895 _propety WITH(NOLOCK)
	LEFT OUTER JOIN dbo._InfoRg15102 _params WITH(NOLOCK)
	ON _propety._IDRRef = _params._Fld15103RRef
	LEFT OUTER JOIN dbo._Reference107 t1 WITH(NOLOCK) -- extra characteristic's values
	ON _params._Fld15106_RRRef = t1._IDRRef
	LEFT OUTER JOIN dbo._Reference123 t2 WITH(NOLOCK) -- brands
	ON _params._Fld15106_RRRef = t2._IDRRef
	LEFT OUTER JOIN dbo._Reference265 t3 WITH(NOLOCK) -- pesrons
	ON _params._Fld15106_RRRef = t3._IDRRef
	LEFT OUTER JOIN dbo._Reference286 t4 WITH(NOLOCK) -- goods
	ON _params._Fld15106_RRRef = t4._IDRRef
	LEFT OUTER JOIN dbo._Reference295 t5 WITH(NOLOCK) -- clinics
	ON _params._Fld15106_RRRef = t5._IDRRef
	LEFT OUTER JOIN dbo._Reference17405 t6 WITH(NOLOCK) -- target for direction clinic
	ON _params._Fld15106_RRRef = t6._IDRRef
	LEFT OUTER JOIN dbo._Reference17871 t7 WITH(NOLOCK) -- strings without size
	ON _params._Fld15106_RRRef = t7._IDRRef
WHERE _propety._Fld18365 <> '' AND _params._Fld15104RRef = @reception

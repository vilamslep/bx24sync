USE dev

DECLARE @clinic BINARY(16)
DECLARE @retail BINARY(16)
DECLARE @optic BINARY(16)
DECLARE @client BINARY(16)

SET @clinic = 0x98AC00025553E65311E3631E1F2C2507
SET @retail = 0x98AC00025553E65311E3630F4F92EDD1
SET @optic = 0x80CA00269E587E4911E3A09A8634E77C
SET @client = ${client}


DECLARE @yearsOffset NUMERIC(4)

SELECT @yearsOffset = offset FROM dbo._YearOffset

------------------------------------------------
--main information
SELECT	
	client._IDRRef AS id,
	CONVERT(VARCHAR(40), client._IDRRef, 2) as originId,
	client._Description AS name,
	client._Fld3808 AS birthday,
	ISNULL(gender._EnumOrder,100) AS gender,
	client._Fld3795 AS isClient,
	client._Fld3797 AS isSuppler,
	client._Fld3801 AS otherRelition,
	client._Fld17995 AS isWorker,
	client._Fld17962 AS isRetireeOrDisabledPerson,
	ISNULL(connectionway._Enumorder, 100) AS connectionway,
	client._Fld18000  AS InRayBanClub,
	CASE 
		WHEN client._IDRRef IN
		(
			SELECT
				_contract._Fld12285RRef
			FROM dbo._Document480 _contract
			WHERE _contract._Posted = 0x01 AND _contract._Fld12285RRef = @client
		) 
		THEN 0x01 
		ELSE 0x00 
	END AS ThereIsContract,
	CASE 
		WHEN client._Fld17241 = 0x01 OR client._Fld17240 = 0x01 OR client._Fld18201 = 0x01 OR client._Fld18200 = 0x01 OR client._Fld18199 = 0x01
		THEN 0x01 
		ELSE 0x00 
	END AS SendAds
INTO #clients
FROM
	dbo._Reference151 AS client WITH(NOLOCK)
	LEFT OUTER JOIN dbo._Enum595 AS gender WITH(NOLOCK)
		ON client._Fld3807RRef = gender._IDRRef
	LEFT OUTER JOIN dbo._Enum18107 connectionway WITH(NOLOCK)
		ON client._Fld18108RRef	= connectionway._IDRRef
WHERE
	client._IDRRef = @client

-------------------------------------------------------------

--- get contacts

DECLARE @typephone BINARY(16)
DECLARE @typeemail BINARY(16)
DECLARE @typeaddress BINARY(16)

SET @typephone = 0xA873CB4AD71D17B2459F9A70D4E2DA66
SET @typeemail = 0x82E6D573EE35D0904BF4D326A84A91D2
SET @typeaddress = 0xAE8167157822C4B643D29FDC57B31A5D

SELECT
	client.id AS clientid,
	contacts._Fld3828 AS phone,
	contacts._Fld3826 AS email
INTO #contacts
FROM
	#clients AS client WITH(NOLOCK)
	LEFT OUTER JOIN dbo._Reference151_VT3817 AS contacts WITH(NOLOCK)
		ON client.id = contacts._Reference151_IDRRef
WHERE
	contacts._Fld3819RRef IN (@typephone, @typeemail)

DECLARE @cursor CURSOR
DECLARE @phonesummary VARCHAR(150)
DECLARE @emailsummary VARCHAR(150)
DECLARE @phone VARCHAR(15)
DECLARE @email VARCHAR(100)

BEGIN 
	SET @cursor = CURSOR FOR 
	SELECT phone, email FROM #contacts
	OPEN @cursor
	FETCH NEXT FROM @cursor INTO @phone, @email
	WHILE @@FETCH_STATUS = 0
	BEGIN
		IF @phone  <> ''
			SET @phonesummary = CONCAT(@phonesummary, @phone, ';')
		IF @email <> ''
			SET @emailsummary = CONCAT(@emailsummary, @email, ';')
		FETCH NEXT FROM @cursor INTO @phone, @email
	END;
END;

DROP TABLE #contacts

------------------------------------------------------------
-- define whose is client: offline, clinic or internet shop
SELECT
	_order._IDRRef AS id,
    _order._Date_Time AS dateObj,
	_client.id AS client,
    _order._Fld17567 AS isInternetClient,
	CASE 
		WHEN _order._Fld17567 = 0x01 
		THEN 0x00 
		WHEN _order._Fld7831RRef = @retail OR _order._Fld7831RRef = @optic 
		THEN 0x01 
		ELSE 0x00 
	END AS isOfflineClient,
	CASE 
		WHEN (_order._Fld7831RRef = @client) 
		THEN 0x01 
		ELSE 0x00 
	END AS isClinicClient
INTO #documents
FROM 
    #clients AS _client WITH(NOLOCK)
    LEFT OUTER JOIN dbo._Document367 AS _order WITH(NOLOCK)
        ON _client.id = _order._Fld7829RRef
WHERE 
	_order._Posted = 0x01

UNION ALL

SELECT
	_reception._IDRRef AS id,
	_reception._Date_Time AS dateObj,
	_client.id AS client,
	0x00 AS isInternetClient,
	CASE 
		WHEN _person._Fld18188RRef  = @retail OR _person._Fld18188RRef  = @optic 
		THEN 0x01 
		ELSE 0x00 
	END,
	CASE 
		WHEN _person._Fld18188RRef = @clinic
		THEN 0x01 
		ELSE 0x00 
	END
FROM #clients AS _client WITH(NOLOCK)
	LEFT OUTER JOIN dev.dbo._Document484 AS _reception WITH(NOLOCK)
		ON _client.id = _reception._Fld12347RRef
	LEFT OUTER JOIN dbo._Reference265 _person WITH(NOLOCK)
		ON (_reception._Fld12348RRef = _person._IDRRef)
WHERE _reception._Posted = 0x01

UNION ALL 

SELECT
	_shipment._IDRRef AS id,
	_shipment._Date_Time AS dateObj,
	_client.id AS client,
	0x00 AS isInternetClient,
	CASE 
		WHEN _shipment._Fld10607RRef = @retail OR _shipment._Fld10607RRef = @optic
		THEN 0x01 
		ELSE 0x00 
	END AS isOfflineClient,
	CASE 
		WHEN _shipment._Fld10607RRef = @clinic 
		THEN 0x01 
		ELSE 0x00 
	END AS isClinicClient
FROM 
	#clients AS _client WITH(NOLOCK)
	LEFT OUTER JOIN dbo._Document431 AS _shipment WITH(NOLOCK)
		ON _client.id = _shipment._Fld10612RRef
WHERE 
	_shipment._Posted = 0x01


SELECT TOP 1
	*
INTO #theFirstDocument
FROM #documents AS document
ORDER BY 
	document.dateObj

DROP TABLE #documents

------------------------------------------------
---define discount value

SELECT
	ISNULL(CAST(CAST(SUM(ammount.expense) AS NUMERIC(33, 8)) AS NUMERIC(27, 2)),0.0) AS expense
INTO #accumSum
FROM (
	SELECT
		clientsSettlement._Period AS Period_,
		clientsSettlement._RecorderTRef AS RecorderTRef,
		clientsSettlement._RecorderRRef AS RecorderRRef,
		ISNULL(CAST(CAST(SUM(
            CASE WHEN clientsSettlement._RecordKind = 0.0 THEN 0.0 ELSE clientsSettlement._Fld16479 END
            ) AS NUMERIC(27, 8)) AS NUMERIC(21, 2)),0.0) AS expense
	
    FROM dbo._AccumRg16475 clientsSettlement WITH(NOLOCK)
		LEFT OUTER JOIN dbo._Reference114 AS analiticKeys WITH(NOLOCK)
			ON (clientsSettlement._Fld16476RRef = analiticKeys._IDRRef)
	WHERE 
		clientsSettlement._Active = 0x01 AND analiticKeys._Fld2988RRef = @client AND clientsSettlement._RecordKind = 1
	
	GROUP BY clientsSettlement._Period,
		clientsSettlement._RecorderTRef,
		clientsSettlement._RecorderRRef
	HAVING (ISNULL( CAST( CAST( SUM(CASE WHEN clientsSettlement._RecordKind = 0.0 THEN 0.0 ELSE clientsSettlement._Fld16479 END ) AS NUMERIC(27, 8) ) AS NUMERIC(21, 2) ),0.0)) <> 0.0 ) AS ammount

DECLARE @ammountSum NUMERIC(15, 2)
DECLARE @discountRayban NUMERIC(2)
DECLARE @discountMedicalThings NUMERIC(2)
DECLARE @discountClinicService NUMERIC(2)

SELECT
	@ammountSum = accum.expense 
FROM #accumSum AS accum 

DROP TABLE #accumSum

SELECT
	@discountMedicalThings  = MAX(T2._Fld5055)
FROM 
	dbo._Reference224_VT5068 Discount WITH(NOLOCK)
		LEFT OUTER JOIN dbo._Reference224 T2 WITH(NOLOCK)
			ON (Discount._Reference224_IDRRef = T2._IDRRef)
		LEFT OUTER JOIN dbo._Reference260 T3 WITH(NOLOCK)
			ON (Discount._Fld5070RRef = T3._IDRRef)
WHERE 
	((T3._Fld5827 <= @ammountSum) 
	AND (T2._Fld17965RRef = 0x911364CC6697B31A4351BE35CD30A94D) 
	AND (T2._Fld5058RRef = 0x87FE01439C28582841B75D998E54D243) 
	AND (T3._Fld5821RRef = 0x8B07C5FB8E31CA8A464CBE1FF6011E7E) 
	AND (T3._Fld5819RRef = 0xA465521B4A08381C45304B8A0E3C631F))
GROUP BY 
	T2._Fld17965RRef

SELECT
	@discountClinicService = MAX(T5._Fld5055)
FROM 
	dbo._Reference224_VT5068 T4 WITH(NOLOCK)
		LEFT OUTER JOIN dbo._Reference224 T5 WITH(NOLOCK)
			ON (T4._Reference224_IDRRef = T5._IDRRef)
		LEFT OUTER JOIN dbo._Reference260 T6 WITH(NOLOCK)
			ON (T4._Fld5070RRef = T6._IDRRef)
WHERE 
	((T6._Fld5827 <= @ammountSum) 
	AND (T5._Fld17965RRef = 0x9D6B19B103D8C976472019B3C44BDD69) 
	AND (T5._Fld5058RRef = 0x87FE01439C28582841B75D998E54D243) 
	AND (T6._Fld5821RRef = 0x8B07C5FB8E31CA8A464CBE1FF6011E7E) 
	AND (T6._Fld5819RRef = 0xA465521B4A08381C45304B8A0E3C631F))
GROUP BY 
	T5._Fld17965RRef

DECLARE @retireDiscontValue NUMERIC(2)
DECLARE @raybanDiscontValue NUMERIC(2)
SET @retireDiscontValue = 10
SET @raybanDiscontValue = 3

------------------------------------------------
-- amount table

SELECT
	LOWER(
		CONCAT(
			SUBSTRING(_client.originId,25,8), 
			'-', SUBSTRING(_client.originId,21,4),
			'-', SUBSTRING(_client.originId,17,4),
			'-', SUBSTRING(_client.originId,1,4), 
			'-', SUBSTRING(_client.originId,5 , 12)
			)
	) AS originId,
	_client.name,
	CASE
		WHEN _client.birthday = '2001-01-01 00:00:00.000'
		THEN NULL
		ELSE FORMAT(DATEADD(YEAR, -@yearsOffset, _client.birthday), 'dd.MM.yyyy', 'ru-RU')
	END AS birthday,
	_client.gender,
	_client.isClient,
	_client.isSuppler,
	_client.otherRelition as otherRelation,
	_client.isRetireeOrDisabledPerson,
	_client.connectionway as connectionWay,
	_client.ThereIsContract AS thereIsContract,
	_client.SendAds AS sendAds,
	_document.isInternetClient,
	_document.isOfflineClient,
	_document.isClinicClient,
	CASE
		WHEN _client.isRetireeOrDisabledPerson = 0x01
		THEN CASE WHEN @discountClinicService > @retireDiscontValue THEN @discountClinicService ELSE @retireDiscontValue END
		ELSE @discountClinicService
	END AS discountClinicService,--кл
	CASE
		WHEN _client.isRetireeOrDisabledPerson = 0x01
		THEN CASE WHEN @discountMedicalThings > @retireDiscontValue THEN @discountMedicalThings ELSE @retireDiscontValue END
		ELSE @discountMedicalThings
	END AS discountMedicalThings,--мо
	CASE 
		WHEN _client.InRayBanClub = 0x01
		THEN @raybanDiscontValue
		ELSE 0
	END AS discountRayban,
	ISNULL(@phonesummary,'') AS phone,
	ISNULL(@email, '') AS email
FROM
	#clients AS _client WITH(NOLOCK)
	LEFT OUTER JOIN #theFirstDocument AS _document WITH(NOLOCK)
		ON _client.id = _document.client

DROP TABLE #clients
DROP TABLE #theFirstDocument


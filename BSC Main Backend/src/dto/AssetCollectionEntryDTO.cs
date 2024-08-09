namespace BSC_Main_Backend.dto;
/// <summary>
/// Local type for use in the AssetCollectionResponseDTO
/// </summary>
/// <param name="Id">Id of GraphicAsset</param>
/// <param name="Transform">Local transform to be applied when asset is displayed as part of collection</param>
public record AssetCollectionEntryDTO(uint Id, TransformDTO Transform);
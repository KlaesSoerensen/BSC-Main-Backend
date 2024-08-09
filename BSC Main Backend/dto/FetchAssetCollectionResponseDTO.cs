namespace BSC_Main_Backend.dto;
using System.Collections.Generic;
/// <summary>
/// Local type for use in the FetchAssetCollectionResponseDTO
/// </summary>
/// <param name="Id">Id of GraphicAsset</param>
/// <param name="Transform">Local transform to be applied when asset is displayed as part of collection</param>
public record AssetCollectionEntryDTO(uint Id, TransformDTO Transform);

/// <param name="Id">Id of AssetCollection</param>
/// <param name="Original">Id of original, non-baked collection</param>
/// <param name="Name">Name of collection, the original is named as "xxxx-original"</param>
/// <param name="Entries">Each GraphicAsset as part of this collection</param>
public record FetchAssetCollectionResponseDTO(uint Id, uint Original, string Name, List<AssetCollectionEntryDTO> Entries);
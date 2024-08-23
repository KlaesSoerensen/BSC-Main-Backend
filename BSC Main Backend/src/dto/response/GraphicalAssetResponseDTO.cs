using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto.response;

/// <param name="Id">Unique identifier of the GraphicalAsset</param>
/// <param name="Alias">Alias for debugging purposes</param>
/// <param name="Type">MIME type of the asset, e.g., "image/png"</param>
/// <param name="Width">Width of the asset (resolution in pixels)</param>
/// <param name="Height">Height of the asset (resolution in pixels)</param>
/// <param name="HasLOD">Indicates whether the asset has multiple Levels of Detail (LOD)</param>
/// <param name="UseCase">Use case of the asset, e.g., "environment", "player", "icon"</param>
/// <param name="Blob">Binary data of the asset, included conditionally based on HasLOD</param>
public record GraphicalAssetResponseDTO(
    uint Id,
    string Alias,
    string Type,  // MIME type
    uint Width,
    uint Height,
    bool HasLOD,
    string UseCase,  // Could map to the AssetUseCase.Value
    Blob Blob  // Optional: Send this only if necessary, or manage its size/serialization
);
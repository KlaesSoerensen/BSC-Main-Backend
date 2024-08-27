using BSC_Main_Backend.Models;
using Microsoft.EntityFrameworkCore;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace BSC_Main_Backend.Data
{
    public class GraphicalAssetRepository : IGraphicalAssetRepository
    {
        private readonly DBContext _context;

        public GraphicalAssetRepository(DBContext context)
        {
            _context = context;
        }

        public async Task<IEnumerable<GraphicalAsset>> GetAllGraphicalAssetsAsync()
        {
            return await _context.GraphicalAssets.ToListAsync();
        }

        public async Task<GraphicalAsset> GetGraphicalAssetByIdAsync(int id)
        {
            return await _context.GraphicalAssets.FindAsync(id);
        }

        public async Task CreateGraphicalAssetAsync(GraphicalAsset graphicalAsset)
        {
            await _context.GraphicalAssets.AddAsync(graphicalAsset);
            await _context.SaveChangesAsync();
        }

        public async Task UpdateGraphicalAssetAsync(GraphicalAsset graphicalAsset)
        {
            _context.GraphicalAssets.Update(graphicalAsset);
            await _context.SaveChangesAsync();
        }

        public async Task DeleteGraphicalAssetAsync(int id)
        {
            var graphicalAsset = await _context.GraphicalAssets.FindAsync(id);
            if (graphicalAsset != null)
            {
                _context.GraphicalAssets.Remove(graphicalAsset);
                await _context.SaveChangesAsync();
            }
        }
    }
}